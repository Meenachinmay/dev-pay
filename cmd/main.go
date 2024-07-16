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
	"sync/atomic"
	"syscall"
	"time"
)

var BATCH_SIZE = 8190

const (
	SmallBatchSize    = 500
	LargeBatchSize    = 8190
	SmallBatchTimeout = 50 * time.Millisecond
	LargeBatchTimeout = 3 * time.Second
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
	requestRate int64
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

	go wp.monitorRequestRate()

	return wp
}

func (wp *CreateAccountWorkerPool) monitorRequestRate() {
	ticker := time.NewTicker(time.Second)
	var count int64
	for range ticker.C {
		atomic.StoreInt64(&wp.requestRate, count)
		count = 0
	}
}

func (wp *CreateAccountWorkerPool) worker(id int) {
	defer wp.wg.Done()
	var batch []Account
	results := make([]chan interface{}, 0)

	processBatch := func() {
		if len(batch) > 0 {
			wp.processBatch(batch, results)
			batch = []Account{}
			results = make([]chan interface{}, 0)
		}
	}

	for {
		select {
		case task, ok := <-wp.tasks:
			if !ok {
				processBatch()
				return
			}
			atomic.AddInt64(&wp.requestRate, 1)
			batch = append(batch, *task.Task)
			results = append(results, task.Result)

			batchSize := SmallBatchSize
			timeout := SmallBatchTimeout
			if atomic.LoadInt64(&wp.requestRate) > int64(wp.workerCount*SmallBatchSize) {
				batchSize = LargeBatchSize
				timeout = LargeBatchTimeout
			}

			if len(batch) >= batchSize {
				processBatch()
			} else {
				time.AfterFunc(timeout, processBatch)
			}
		}
	}
}

func (wp *CreateAccountWorkerPool) processBatch(batch []Account, results []chan interface{}) {
	client := wp.clientPool.GetClient()

	_, err := client.CreateAccounts(batch)
	if err != nil {
		log.Printf("Batch error: %v", err)
		for _, resultChan := range results {
			resultChan <- fmt.Errorf("failed to create accounts due to TB client error: %v", err)
			close(resultChan) // Ensure channels are closed after writing to prevent leaks
		}
		return
	}

	for _, resultChan := range results {
		resultChan <- "Account created successfully"
		close(resultChan) // Ensure channels are closed after writing to prevent leaks
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
	accounts       []Account
	accountResults map[string]chan string // Map to store results channels by batch
	mu             sync.Mutex
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
			return nil, status.Errorf(codes.ResourceExhausted, "server is overloaded, please try again later")
		}

		select {
		case result := <-task.Result:
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
		log.Printf("Create account error: %v", err)
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
		MaxRequests: 8190,
		Interval:    3 * time.Second,
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
		accounts:       make([]Account, 0),
		accountResults: make(map[string]chan string),
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
