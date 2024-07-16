package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	payments "dev-pay/grpc-proto"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	. "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	BatchSize    = 8190
	QueueSize    = 100000
	NumWorkers   = 14
	OpTimeout    = 5 * time.Second
	MinBatchSize = 4000
)

type AccountJob struct {
	Account Account
}

type createAccountServer struct {
	payments.UnimplementedCreateAccountServiceServer
	client   tigerbeetle_go.Client
	jobQueue chan AccountJob
	workerWg sync.WaitGroup
	quit     chan struct{}
}

func (s *createAccountServer) CreateAccount(ctx context.Context, req *payments.CreateAccountRequest) (*payments.CreateAccountResponse, error) {
	tbAccount := convertGrpcAccountToTigerBeetleAccount(req.Account)

	job := AccountJob{
		Account: tbAccount,
	}

	select {
	case s.jobQueue <- job:
		return &payments.CreateAccountResponse{Results: "Account queued for creation"}, nil
	default:
		log.Println("Server overloaded, rejecting request")
		return nil, status.Errorf(codes.ResourceExhausted, "server is overloaded, please try again later")
	}
}

func (s *createAccountServer) worker(id int) {
	defer s.workerWg.Done()
	log.Printf("Worker %d started", id)

	batch := make([]Account, 0, BatchSize)
	adaptiveMinBatchSize := MinBatchSize
	tickerInterval := 2 * time.Second

	flushBatch := func() {
		batchSize := len(batch)
		if batchSize >= adaptiveMinBatchSize {
			log.Printf("Worker %d flushing batch of size %d", id, batchSize)
			_, cancel := context.WithTimeout(context.Background(), OpTimeout)
			defer cancel()

			start := time.Now()
			_, err := s.client.CreateAccounts(batch)
			duration := time.Since(start)

			if err != nil {
				log.Printf("Worker %d failed to create accounts: %v", id, err)
			} else {
				// Adjust adaptive parameters based on performance
				if duration < 500*time.Millisecond && adaptiveMinBatchSize < BatchSize {
					adaptiveMinBatchSize = min(adaptiveMinBatchSize+500, BatchSize)
					tickerInterval = min(tickerInterval+500*time.Millisecond, 5*time.Second)
				} else if duration > 1*time.Second && adaptiveMinBatchSize > MinBatchSize {
					adaptiveMinBatchSize = max(adaptiveMinBatchSize-500, MinBatchSize)
					tickerInterval = max(tickerInterval-500*time.Millisecond, 1*time.Second)
				}
			}

			batch = batch[:0]
		}
	}

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case job, ok := <-s.jobQueue:
			if !ok {
				if len(batch) > 0 {
					_, err := s.client.CreateAccounts(batch)
					if err != nil {
						log.Printf("Worker %d failed to create final batch: %v", id, err)
					}
				}
				log.Printf("Worker %d shutting down", id)
				return
			}
			batch = append(batch, job.Account)
			if len(batch) >= BatchSize {
				flushBatch()
				ticker.Reset(tickerInterval)
			}
		case <-ticker.C:
			flushBatch()
			ticker.Reset(tickerInterval)
		case <-s.quit:
			if len(batch) > 0 {
				_, err := s.client.CreateAccounts(batch)
				if err != nil {
					log.Printf("Worker %d failed to create final batch: %v", id, err)
				}
			}
			log.Printf("Worker %d received quit signal", id)
			return
		}
	}
}

func convertGrpcAccountToTigerBeetleAccount(grpcAcc *payments.Account) Account {
	return Account{
		ID:     ToUint128(grpcAcc.Id),
		Ledger: grpcAcc.Ledger,
		Code:   uint16(grpcAcc.Code),
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := os.Getenv("TB_ADDRESS")
	if port == "" {
		port = "3000"
	}

	client, err := tigerbeetle_go.NewClient(ToUint128(uint64(0)), []string{port}, BatchSize)
	if err != nil {
		log.Fatalf("Error creating TigerBeetle client: %s", err)
	}
	defer client.Close()

	log.Println("Connected to TigerBeetle")

	server := &createAccountServer{
		client:   client,
		jobQueue: make(chan AccountJob, QueueSize),
		quit:     make(chan struct{}),
	}

	for i := 0; i < NumWorkers; i++ {
		server.workerWg.Add(1)
		go server.worker(i)
	}

	grpcServerOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Hour,
			Timeout:           20 * time.Second,
		}),
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	payments.RegisterCreateAccountServiceServer(grpcServer, server)
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

	close(server.quit)
	grpcServer.GracefulStop()
	close(server.jobQueue)
	server.workerWg.Wait()
	log.Println("Server gracefully stopped")
}
