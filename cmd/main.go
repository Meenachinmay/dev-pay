package main

import (
	"context"
	payments "dev-pay/grpc-proto"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sony/gobreaker"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	. "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type CreateAccountTask struct {
	Task   *Account
	Result chan interface{}
	Ctx    context.Context
}

type TigerBeetleClientPool struct {
	clients []tigerbeetle_go.Client
	mu      sync.Mutex
	index   int
}

func NewTigerBeetleClientPool(size int, address string) (*TigerBeetleClientPool, error) {
	pool := &TigerBeetleClientPool{
		clients: make([]tigerbeetle_go.Client, size),
	}

	for i := 0; i < size; i++ {
		client, err := tigerbeetle_go.NewClient(ToUint128(uint64(0)), []string{address}, 8192)
		if err != nil {
			return nil, fmt.Errorf("failed to create TigerBeetle client %d: %v", i, err)
		}
		pool.clients[i] = client
	}

	return pool, nil
}

func (p *TigerBeetleClientPool) GetClient() tigerbeetle_go.Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	client := p.clients[p.index]
	p.index = (p.index + 1) % len(p.clients)
	return client
}

func (p *TigerBeetleClientPool) Close() {
	for _, client := range p.clients {
		client.Close()
	}
}

type CreateAccountWorkerPool struct {
	tasks       chan CreateAccountTask
	clientPool  *TigerBeetleClientPool
	workerCount int
	wg          sync.WaitGroup
}

func NewCreateAccountWorkerPool(workerCount int, taskQueueSize int, clientPool *TigerBeetleClientPool) *CreateAccountWorkerPool {
	wp := &CreateAccountWorkerPool{
		tasks:       make(chan CreateAccountTask, taskQueueSize),
		clientPool:  clientPool,
		workerCount: workerCount,
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	return wp
}

func (wp *CreateAccountWorkerPool) worker(id int) {
	defer wp.wg.Done()
	for task := range wp.tasks {
		client := wp.clientPool.GetClient()
		select {
		case <-task.Ctx.Done():
			task.Result <- fmt.Errorf("task cancelled")
		default:
			res, err := client.CreateAccounts([]Account{*task.Task})
			if err != nil {
				task.Result <- fmt.Errorf("failed to create account due to TB client error: %v", err)
				continue
			}

			if len(res) > 0 {
				task.Result <- fmt.Errorf("error creating account %d: %s", res[0].Index, res[0].Result)
			} else {
				task.Result <- "Account created successfully"
			}
		}
	}
}

func (wp *CreateAccountWorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}

func convertGrpcAccountToTigerBeetleAccount(grpcAcc *payments.Account) *Account {
	return &Account{
		ID:     ToUint128(grpcAcc.Id),
		Ledger: grpcAcc.Ledger,
		Code:   uint16(grpcAcc.Code),
	}
}

type createAccountServer struct {
	payments.UnimplementedCreateAccountServiceServer
	accountPool    *CreateAccountWorkerPool
	circuitBreaker *gobreaker.CircuitBreaker
}

func (s *createAccountServer) CreateAccount(ctx context.Context, req *payments.CreateAccountRequest) (*payments.CreateAccountResponse, error) {
	result, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		tbAccount := convertGrpcAccountToTigerBeetleAccount(req.Account)

		task := CreateAccountTask{
			Task:   tbAccount,
			Result: make(chan interface{}, 1),
			Ctx:    ctx,
		}

		select {
		case s.accountPool.tasks <- task:
			// Task submitted successfully
		default:
			// Worker pool is full, implement backpressure
			return nil, status.Errorf(codes.ResourceExhausted, "server is overloaded, please try again later")
		}

		select {
		case result := <-task.Result:
			close(task.Result)
			switch v := result.(type) {
			case error:
				return nil, status.Errorf(codes.Internal, "failed to create account: %v", v)
			case string:
				return &payments.CreateAccountResponse{Results: v}, nil
			default:
				return nil, status.Errorf(codes.Internal, "unexpected result type")
			}
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "request timed out")
		}
	})

	if err != nil {
		return nil, err
	}

	response, ok := result.(*payments.CreateAccountResponse)
	if !ok {
		return nil, status.Errorf(codes.Internal, "unexpected response type")
	}

	return response, nil
}

func main() {
	// Set GOMAXPROCS to the number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "CreateAccount Circuit Breaker",
		MaxRequests: 100,
		Interval:    5 * time.Second,
		Timeout:     3 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})

	// Connect to TigerBeetle
	port := os.Getenv("TB_ADDRESS")
	if port == "" {
		port = "3000"
	}

	clientPool, err := NewTigerBeetleClientPool(runtime.NumCPU(), port)
	if err != nil {
		log.Fatalf("Error creating TigerBeetle client pool: %s", err)
	}
	defer clientPool.Close()

	log.Println("Connected to TigerBeetle")

	// Create worker pool
	workerCount := runtime.NumCPU() * 2 // Adjust this based on your needs
	taskQueueSize := 10000              // Adjust this based on your needs
	accountPool := NewCreateAccountWorkerPool(workerCount, taskQueueSize, clientPool)
	defer accountPool.Shutdown()

	grpcServerOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Hour,
			Timeout:           20 * time.Second,
		}),
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	payments.RegisterCreateAccountServiceServer(grpcServer, &createAccountServer{
		accountPool:    accountPool,
		circuitBreaker: cb,
	})
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("gRPC server listening on port 50051")
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	grpcServer.GracefulStop()
	log.Println("Server gracefully stopped")
}
