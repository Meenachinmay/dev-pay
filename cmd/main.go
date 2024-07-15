package main

import (
	"context"
	payments "dev-pay/grpc-proto"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/sony/gobreaker"
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
	Task interface{}
}

type LookUpAccountTask struct {
	Task interface{}
}

type CreateTransferTask struct {
	Task interface{}
}

type CreateAccountWorkerPool struct {
	tasks chan CreateAccountTask
	wg    sync.WaitGroup
}

type LookUpAccountWorkerPool struct {
	tasks chan LookUpAccountTask
	wg    sync.WaitGroup
}

type CreateTransferWorkerPool struct {
	tasks chan CreateTransferTask
	wg    sync.WaitGroup
}

func NewCreateAccountWorkerPool(workerCount int, taskQueueSize int) *CreateAccountWorkerPool {
	wp := &CreateAccountWorkerPool{
		tasks: make(chan CreateAccountTask, taskQueueSize),
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
func NewCreateTransferWorkerPool(workerCount int, taskQueueSize int) *CreateTransferWorkerPool {
	wp := &CreateTransferWorkerPool{
		tasks: make(chan CreateTransferTask, taskQueueSize),
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

func NewLookUpAccountWorkerPool(workerCount int, taskQueueSize int) *LookUpAccountWorkerPool {
	wp := &LookUpAccountWorkerPool{
		tasks: make(chan LookUpAccountTask, taskQueueSize),
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

func (s *createTransferServer) CreateAccounts(ctx context.Context, req *payments.CreateAccountsRequest) (*payments.CreateAccountsResponse, error) {
	resultChan := make(chan string)
	payload, _ := json.Marshal(req)

}

func (s *lookUpAccountServer) LookupAccounts(ctx context.Context, req *payments.LookupAccountsRequest) (*payments.LookupAccountsResponse, error) {
	resultChan := make(chan []Account)
	payload, _ := json.Marshal(req)

}

func (s *createTransferServer) CreateTransfers(ctx context.Context, req *payments.CreateTransfersRequest) (*payments.CreateTransfersResponse, error) {
	resultChan := make(chan string)
	payload, _ := json.Marshal(req)

}

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

	// connect to database

	// Create separate worker pools for each service
	accountPool := NewCreateAccountWorkerPool(runtime.NumCPU(), 8190)
	lookupPool := NewLookUpAccountWorkerPool(runtime.NumCPU(), 8190)
	transferPool := NewCreateTransferWorkerPool(runtime.NumCPU(), 8190)
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
