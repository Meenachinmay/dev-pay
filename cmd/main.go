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
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
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
}

type LookUpAccountTask struct {
	Task   []Account
	Result chan interface{}
}

type CreateTransferTask struct {
	Task   []Transfer
	Result chan interface{}
}

type CreateAccountWorkerPool struct {
	tasks    chan CreateAccountTask
	wg       sync.WaitGroup
	tbClient tigerbeetle_go.Client
}

type LookUpAccountWorkerPool struct {
	tasks    chan LookUpAccountTask
	wg       sync.WaitGroup
	tbClient tigerbeetle_go.Client
}

type CreateTransferWorkerPool struct {
	tasks    chan CreateTransferTask
	wg       sync.WaitGroup
	tbClient tigerbeetle_go.Client
}

func NewCreateAccountWorkerPool(workerCount int, taskQueueSize int, client tigerbeetle_go.Client) *CreateAccountWorkerPool {
	wp := &CreateAccountWorkerPool{
		tasks:    make(chan CreateAccountTask, taskQueueSize),
		tbClient: client,
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for task := range wp.tasks {
				res, err := wp.tbClient.CreateAccounts([]Account{*task.Task})
				if err != nil {
					task.Result <- fmt.Errorf("failed to create account due to TB client error: %v", err)
					continue
				}

				for _, err := range res {
					task.Result <- fmt.Errorf("error creating account %d: %s", err.Index, err.Result)
				}

				task.Result <- "Account created successfully"
			}
		}(i)
	}
	return wp
}

func NewCreateTransferWorkerPool(workerCount int, taskQueueSize int, client tigerbeetle_go.Client) *CreateTransferWorkerPool {
	wp := &CreateTransferWorkerPool{
		tasks:    make(chan CreateTransferTask, taskQueueSize),
		tbClient: client,
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go func(id int) {
			// logic
			wp.wg.Done()
		}(i)
	}
	return wp
}

func NewLookUpAccountWorkerPool(workerCount int, taskQueueSize int, client tigerbeetle_go.Client) *LookUpAccountWorkerPool {
	wp := &LookUpAccountWorkerPool{
		tasks:    make(chan LookUpAccountTask, taskQueueSize),
		tbClient: client,
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go func(id int) {
			// logic
			wp.wg.Done()
		}(i)
	}
	return wp
}

func (wp *CreateAccountWorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}
func (wp *CreateTransferWorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}
func (wp *LookUpAccountWorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}

// Helper function to convert a gRPC Account to a TigerBeetle Account
func convertGrpcAccountToTigerBeetleAccount(grpcAcc *payments.Account) *Account {
	return &Account{
		ID:     ToUint128(grpcAcc.Id),
		Ledger: grpcAcc.Ledger,
		Code:   uint16(grpcAcc.Code), // Assuming Code is a uint16 in TigerBeetle definition
	}
}

func (s *createAccountServer) CreateAccount(ctx context.Context, req *payments.CreateAccountRequest) (*payments.CreateAccountResponse, error) {
	resultChan := make(chan interface{})
	tbAccount := convertGrpcAccountToTigerBeetleAccount(req.Account)

	// prepare the request
	task := &CreateAccountTask{
		Task:   tbAccount,
		Result: resultChan,
	}

	// submit the task to tasks channel along with result channel
	s.accountPool.SubmitTask(*task)

	// listen for the result after from createAccount worker pool
	select {
	case response := <-resultChan:
		log.Printf("CreateAccount response: %v", response)
		close(resultChan)
		return &payments.CreateAccountResponse{Results: response.(string)}, nil
	case <-time.After(time.Second * 5):
		close(resultChan)
		return nil, fmt.Errorf("CreateAccount timeout")
	}

}

//func (s *lookUpAccountServer) LookupAccounts(ctx context.Context, req *payments.LookupAccountsRequest) (*payments.LookupAccountsResponse, error) {
//	resultChan := make(chan []Account)
//	payload, _ := json.Marshal(req)
//
//}
//
//func (s *createTransferServer) CreateTransfers(ctx context.Context, req *payments.CreateTransfersRequest) (*payments.CreateTransfersResponse, error) {
//	resultChan := make(chan string)
//	payload, _ := json.Marshal(req)
//
//}

type createAccountServer struct {
	payments.UnimplementedCreateAccountServiceServer
	accountPool    *CreateAccountWorkerPool
	circuitBreaker *gobreaker.CircuitBreaker
}

type createTransferServer struct {
	payments.UnimplementedCreateTransferServiceServer
	transferPool *CreateTransferWorkerPool
}

type lookUpAccountServer struct {
	payments.UnimplementedTransactionsLookUpServiceServer
	lookupPool *LookUpAccountWorkerPool
}

func (wp *CreateAccountWorkerPool) SubmitTask(task CreateAccountTask) {
	wp.tasks <- task
}
func (wp *CreateTransferWorkerPool) SubmitTask(task CreateTransferTask) {
	wp.tasks <- task
}
func (wp *LookUpAccountWorkerPool) SubmitTask(task LookUpAccountTask) {
	wp.tasks <- task
}

func main() {
	// Set GOMAXPROCS to the number of available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a circuit breaker
	_ = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "gRPC Circuit Breaker",
		MaxRequests: 1000,             // Allow 1000 requests during the half-open state
		Interval:    30 * time.Second, // Reset the failure count after this interval
		Timeout:     10 * time.Second, // Timeout duration for the open state
	})

	// connect to Tiger-beetle here
	port := os.Getenv("TB_ADDRESS")
	if port == "" {
		port = "3000"
	}

	client, err := tigerbeetle_go.NewClient(ToUint128(0), []string{port}, 8192)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	} else {
		log.Println("Connected to Tigerbeetle")
	}
	defer client.Close()

	// Create separate worker pools for each service
	accountPool := NewCreateAccountWorkerPool(runtime.NumCPU(), 8190, client)
	lookupPool := NewLookUpAccountWorkerPool(runtime.NumCPU(), 8190, client)
	transferPool := NewCreateTransferWorkerPool(runtime.NumCPU(), 8190, client)
	defer accountPool.Shutdown()
	defer lookupPool.Shutdown()
	defer transferPool.Shutdown()

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
		accountPool: accountPool,
	})
	payments.RegisterCreateTransferServiceServer(grpcServer, &createTransferServer{
		transferPool: transferPool,
	})
	payments.RegisterTransactionsLookUpServiceServer(grpcServer, &lookUpAccountServer{
		lookupPool: lookupPool,
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
