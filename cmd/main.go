package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
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
	QueueSize    = 1000000
	NumWorkers   = 14
	OpTimeout    = 5 * time.Second
	MinBatchSize = 4000
)

type AccountJob struct {
	Account Account
}

type createAccountServer struct {
	payments.UnimplementedCreateAccountServiceServer
	payments.UnimplementedTransactionsLookUpServiceServer
	client        tigerbeetle_go.Client
	jobQueue      chan AccountJob
	workerWg      sync.WaitGroup
	quit          chan struct{}
	processedJobs int64
}

type TransferJob struct {
	Transfer Transfer
}

type createTransferServer struct {
	payments.UnimplementedCreateTransferServiceServer
	client        tigerbeetle_go.Client
	jobQueue      chan TransferJob
	workerWg      sync.WaitGroup
	quit          chan struct{}
	processedJobs int64
}

func (s *createTransferServer) CreateTransfers(ctx context.Context, req *payments.CreateTransfersRequest) (*payments.CreateTransfersResponse, error) {
	results := make([]string, len(req.Transfers))
	for i, transfer := range req.Transfers {
		tbTransfer := convertGrpcTransferToTigerBeetleTransfer(transfer)

		select {
		case s.jobQueue <- TransferJob{Transfer: tbTransfer}:
			results[i] = "Transfer queued for creation"
		default:
			// Queue is full, wait a bit and try again
			time.Sleep(10 * time.Millisecond)
			select {
			case s.jobQueue <- TransferJob{Transfer: tbTransfer}:
				results[i] = "Transfer queued for creation after delay"
			default:
				results[i] = "Server overloaded, transfer creation failed"
			}
		}
	}
	return &payments.CreateTransfersResponse{Results: results}, nil
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

func (s *createAccountServer) CreateAccountBatch(ctx context.Context, req *payments.CreateAccountBatchRequest) (*payments.CreateAccountBatchResponse, error) {
	results := make([]string, len(req.Accounts))
	for i, account := range req.Accounts {
		tbAccount := convertGrpcAccountToTigerBeetleAccount(account)

		select {
		case s.jobQueue <- AccountJob{Account: tbAccount}:
			results[i] = "Account queued for creation"
		default:
			// Queue is full, wait a bit and try again
			time.Sleep(10 * time.Millisecond)
			select {
			case s.jobQueue <- AccountJob{Account: tbAccount}:
				results[i] = "Account queued for creation after delay"
			default:
				results[i] = "Server overloaded, account creation failed"
			}
		}
	}
	return &payments.CreateAccountBatchResponse{Results: results}, nil
}

func (s *createAccountServer) LookupAccounts(ctx context.Context, req *payments.LookupAccountsRequest) (*payments.LookupAccountsResponse, error) {
	tbAccountIDs := make([]Uint128, len(req.AccountIds))
	for i, id := range req.AccountIds {
		tbAccountIDs[i] = ToUint128(id)
	}

	accounts, err := s.client.LookupAccounts(tbAccountIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to lookup accounts: %v", err)
	}

	response := &payments.LookupAccountsResponse{
		Accounts: make([]*payments.AccountResponse, len(accounts)),
	}

	for i, account := range accounts {
		response.Accounts[i] = &payments.AccountResponse{
			Id:             convertUint128ToProto(account.ID),
			DebitsPending:  convertUint128ToProto(account.DebitsPending),
			DebitsPosted:   convertUint128ToProto(account.DebitsPosted),
			CreditsPending: convertUint128ToProto(account.CreditsPending),
			CreditsPosted:  convertUint128ToProto(account.CreditsPosted),
			UserData128:    convertUint128ToProto(account.UserData128),
			UserData64:     account.UserData64,
			UserData32:     account.UserData32,
			Reserved:       account.Reserved,
			Ledger:         account.Ledger,
			Code:           uint32(account.Code),
			Flags:          uint32(account.Flags),
			Timestamp:      account.Timestamp,
		}
	}

	return response, nil
}

func convertUint128ToProto(value Uint128) *payments.Uint128 {
	bytes := value.Bytes()
	return &payments.Uint128{
		High: binary.LittleEndian.Uint64(bytes[8:]),
		Low:  binary.LittleEndian.Uint64(bytes[:8]),
	}
}

func (s *createAccountServer) worker(id int) {
	defer s.workerWg.Done()
	log.Printf("Worker %d started", id)

	batch := make([]Account, 0, BatchSize)
	adaptiveMinBatchSize := MinBatchSize
	tickerInterval := 500 * time.Millisecond

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

				// IF error, reduce batch size to process smaller chunks
				adaptiveMinBatchSize = max(adaptiveMinBatchSize/2, MinBatchSize)
			} else {
				atomic.AddInt64(&s.processedJobs, int64(batchSize))
				// Adjust adaptive parameters based on performance
				if duration < 200*time.Millisecond && adaptiveMinBatchSize < BatchSize {
					adaptiveMinBatchSize = min(adaptiveMinBatchSize+1000, BatchSize)
					tickerInterval = min(tickerInterval+100*time.Millisecond, 2*time.Second)
				} else if duration > 500*time.Millisecond && adaptiveMinBatchSize > MinBatchSize {
					adaptiveMinBatchSize = max(adaptiveMinBatchSize-1000, MinBatchSize)
					tickerInterval = max(tickerInterval-100*time.Millisecond, 100*time.Millisecond)
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
					_, cancel := context.WithTimeout(context.Background(), OpTimeout)
					_, err := s.client.CreateAccounts(batch)
					cancel()
					if err != nil {
						log.Printf("Worker %d failed to create final batch: %v", id, err)
					} else {
						atomic.AddInt64(&s.processedJobs, int64(len(batch)))
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
				_, cancel := context.WithTimeout(context.Background(), OpTimeout)
				_, err := s.client.CreateAccounts(batch)
				cancel()
				if err != nil {
					log.Printf("Worker %d failed to create final batch: %v", id, err)
				} else {
					atomic.AddInt64(&s.processedJobs, int64(len(batch)))
				}
			}
			log.Printf("Worker %d received quit signal", id)
			return
		}
	}
}

func (s *createTransferServer) worker(id int) {
	defer s.workerWg.Done()
	log.Printf("Worker %d started", id)

	batch := make([]Transfer, 0, BatchSize)
	adaptiveMinBatchSize := MinBatchSize
	tickerInterval := 500 * time.Millisecond

	flushBatch := func() {
		batchSize := len(batch)
		if batchSize >= adaptiveMinBatchSize {
			log.Printf("Worker %d flushing batch of size %d", id, batchSize)
			_, cancel := context.WithTimeout(context.Background(), OpTimeout)
			defer cancel()

			start := time.Now()
			_, err := s.client.CreateTransfers(batch)
			duration := time.Since(start)

			if err != nil {
				log.Printf("Worker %d failed to create accounts: %v", id, err)

				// IF error, reduce batch size to process smaller chunks
				adaptiveMinBatchSize = max(adaptiveMinBatchSize/2, MinBatchSize)
			} else {
				atomic.AddInt64(&s.processedJobs, int64(batchSize))
				// Adjust adaptive parameters based on performance
				if duration < 200*time.Millisecond && adaptiveMinBatchSize < BatchSize {
					adaptiveMinBatchSize = min(adaptiveMinBatchSize+1000, BatchSize)
					tickerInterval = min(tickerInterval+100*time.Millisecond, 2*time.Second)
				} else if duration > 500*time.Millisecond && adaptiveMinBatchSize > MinBatchSize {
					adaptiveMinBatchSize = max(adaptiveMinBatchSize-1000, MinBatchSize)
					tickerInterval = max(tickerInterval-100*time.Millisecond, 100*time.Millisecond)
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
					_, cancel := context.WithTimeout(context.Background(), OpTimeout)
					_, err := s.client.CreateTransfers(batch)
					cancel()
					if err != nil {
						log.Printf("Worker %d failed to create final batch: %v", id, err)
					} else {
						atomic.AddInt64(&s.processedJobs, int64(len(batch)))
					}
				}
				log.Printf("Worker %d shutting down", id)
				return
			}
			batch = append(batch, job.Transfer)
			if len(batch) >= BatchSize {
				flushBatch()
				ticker.Reset(tickerInterval)
			}
		case <-ticker.C:
			flushBatch()
			ticker.Reset(tickerInterval)
		case <-s.quit:
			if len(batch) > 0 {
				_, cancel := context.WithTimeout(context.Background(), OpTimeout)
				_, err := s.client.CreateTransfers(batch)
				cancel()
				if err != nil {
					log.Printf("Worker %d failed to create final batch: %v", id, err)
				} else {
					atomic.AddInt64(&s.processedJobs, int64(len(batch)))
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

func convertGrpcTransferToTigerBeetleTransfer(grpcTransfer *payments.Transfer) Transfer {
	return Transfer{
		ID:              ToUint128(grpcTransfer.Id),
		DebitAccountID:  ToUint128(grpcTransfer.DebitAccountId),
		CreditAccountID: ToUint128(grpcTransfer.CreditAccountId),
		Amount:          ToUint128(grpcTransfer.Amount),
		Ledger:          grpcTransfer.Ledger,
		Code:            uint16(grpcTransfer.Code),
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

	// Create two accounts
	res, err := client.CreateAccounts([]Account{
		{
			ID:     ToUint128(1),
			Ledger: 1,
			Code:   1,
			Flags:  AccountFlags{}.ToUint16(), // No restrictions
		},
		{
			ID:     ToUint128(2),
			Ledger: 1,
			Code:   1,
			Flags:  AccountFlags{}.ToUint16(), // No restrictions
		},
		{
			ID:     ToUint128(3),
			Ledger: 1,
			Code:   1,
			Flags:  AccountFlags{}.ToUint16(), // No restrictions
		},
		{
			ID:     ToUint128(4),
			Ledger: 1,
			Code:   1,
			Flags:  AccountFlags{}.ToUint16(), // No restrictions
		},
		{
			ID:     ToUint128(7000000),
			Ledger: 1,
			Code:   1,
			Flags:  AccountFlags{}.ToUint16(), // No restrictions
		},
		{
			ID:     ToUint128(7000001),
			Ledger: 1,
			Code:   1,
			Flags:  AccountFlags{}.ToUint16(), // No restrictions
		},
	})
	if err != nil {
		log.Fatalf("Error creating accounts: %s", err)
	} else {
		log.Println("Test Accounts created")
	}

	for _, err := range res {
		log.Fatalf("Error creating account %d: %s", err.Index, err.Result)
	}

	// Check the sums for both accounts
	_accounts, err := client.LookupAccounts([]Uint128{ToUint128(1), ToUint128(2), ToUint128(3), ToUint128(4)})
	if err != nil {
		log.Fatalf("Could not fetch accounts: %s", err)
	}
	assert(len(_accounts) == 4, "accounts")

	transferRes, err := client.CreateTransfers([]Transfer{
		{
			ID:              ToUint128(7000000),
			DebitAccountID:  ToUint128(7000000),
			CreditAccountID: ToUint128(1),
			Amount:          ToUint128(100000000),
			Ledger:          1,
			Code:            1,
		},
		{
			ID:              ToUint128(7000001),
			DebitAccountID:  ToUint128(7000001),
			CreditAccountID: ToUint128(2),
			Amount:          ToUint128(100000000),
			Ledger:          1,
			Code:            1,
		},
	})
	if err != nil {
		log.Fatalf("Error creating transfer: %s", err)
	}

	for _, err := range transferRes {
		log.Fatalf("Error creating transfer: %s", err.Result)
	}

	// Check the sums for both accounts
	accounts, err := client.LookupAccounts([]Uint128{ToUint128(1), ToUint128(2)})
	if err != nil {
		log.Fatalf("Could not fetch accounts: %s", err)
	}
	assert(len(accounts) == 2, "accounts")

	fmt.Printf("Test Accounts status before load start %+v\n\n", accounts)

	server := &createAccountServer{
		client:   client,
		jobQueue: make(chan AccountJob, QueueSize),
		quit:     make(chan struct{}),
	}

	for i := 0; i < NumWorkers; i++ {
		server.workerWg.Add(1)
		go server.worker(i)
	}

	transferServer := &createTransferServer{
		client:   client,
		jobQueue: make(chan TransferJob, QueueSize),
		quit:     make(chan struct{}),
	}

	for i := 0; i < NumWorkers; i++ {
		transferServer.workerWg.Add(1)
		go transferServer.worker(i)
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
	payments.RegisterCreateTransferServiceServer(grpcServer, transferServer)
	payments.RegisterTransactionsLookUpServiceServer(grpcServer, server)
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

	go server.monitorProgress()
	go transferServer.monitorProgress()
	go verifyAccountBalances(client)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	close(server.quit)
	grpcServer.GracefulStop()
	close(server.jobQueue)
	server.workerWg.Wait()
	log.Printf("Server gracefully stopped. Total accounts processed: %d", atomic.LoadInt64(&server.processedJobs))

}

// monitor for create account service
func (s *createAccountServer) monitorProgress() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var lastProcessed int64
	lastTime := time.Now()

	for {
		select {
		case <-ticker.C:
			currentProcessed := atomic.LoadInt64(&s.processedJobs)
			currentTime := time.Now()
			duration := currentTime.Sub(lastTime)
			rate := float64(currentProcessed-lastProcessed) / duration.Seconds()

			log.Printf("Processed: %d, Queue size: %d, Processing rate: %.2f accounts/second",
				currentProcessed, len(s.jobQueue), rate)

			lastProcessed = currentProcessed
			lastTime = currentTime
		case <-s.quit:
			return
		}
	}
}

// monitor for create transfer service
func (s *createTransferServer) monitorProgress() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var lastProcessed int64
	lastTime := time.Now()

	for {
		select {
		case <-ticker.C:
			currentProcessed := atomic.LoadInt64(&s.processedJobs)
			currentTime := time.Now()
			duration := currentTime.Sub(lastTime)
			rate := float64(currentProcessed-lastProcessed) / duration.Seconds()

			log.Printf("Processed: %d, Queue size: %d, Processing rate: %.2f accounts/second",
				currentProcessed, len(s.jobQueue), rate)

			lastProcessed = currentProcessed
			lastTime = currentTime
		case <-s.quit:
			return
		}
	}
}

func verifyAccountBalances(client tigerbeetle_go.Client) {
	accountIDs := []Uint128{
		ToUint128(1),
		ToUint128(2),
		ToUint128(3),
		ToUint128(4),
	}

	accounts, err := client.LookupAccounts(accountIDs)
	if err != nil {
		log.Printf("Failed to lookup accounts: %v", err)
		return
	}

	for _, account := range accounts {
		log.Printf("Account %v: Credits Posted = %v, Debits Posted = %v",
			account.ID, account.CreditsPosted, account.DebitsPosted)
	}
}

func assert(condition bool, message string) {
	if !condition {
		log.Fatalf("Assertion failed: %s", message)
	}
}
