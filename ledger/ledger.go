package ledger

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Monrevil/simplified-market-ledger/api"
	"github.com/Monrevil/simplified-market-ledger/ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/ledger/repository/postgres"
	"github.com/Monrevil/simplified-market-ledger/ledger/transactions"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Ledger struct {
	r   *postgres.PostgresRepository
	log *zap.SugaredLogger
	api.UnimplementedLedgerServer
}

func Serve(address string) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	zapLogger := logger.Sugar()
	zapLogger.Infof("Starting a gRPC server on %s...", address)

	conn, err := net.Listen("tcp", address)
	if err != nil {
		zapLogger.Fatalf("Cannot listen on address %s", address)
	}
	s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		loggingMiddleware(logger),
		grpc_zap.UnaryServerInterceptor(logger),
	),
	)
	ledgerServer := &Ledger{
		r:   postgres.NewPostgresRepository(),
		log: zapLogger,
	}
	api.RegisterLedgerServer(s, ledgerServer)

	// Perform graceful shutdown if interrupted.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		stop := <-sigCh
		ledgerServer.log.Infof("Got %v signal, attempting graceful shutdown", stop)
		s.GracefulStop()
		ledgerServer.log.Info("Shut down gRPC server successfully")
	}()

	ledgerServer.log.Infof("gRPC server is listening and serving on %v", address)
	if err := s.Serve(conn); err != nil {
		ledgerServer.log.Fatalf("failed to serve: %v", err)
	}
}

func loggingMiddleware(logger *zap.Logger) grpc.UnaryServerInterceptor {
	loggingInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Info("received a gRPC request",
			zap.String("handler", info.FullMethod),
			zap.String("request", fmt.Sprintf("%+v\n", req)),
		)
		resp, err := handler(ctx, req)
		logger.Info("sending a gRPC response",
			zap.String("handler", info.FullMethod),
			zap.String("response", fmt.Sprintf("%+v", resp)),
		)
		return resp, err
	}
	return loggingInterceptor
}

func (l *Ledger) NewIssuer(ctx context.Context, req *api.NewIssuerReq) (*api.NewIssuerResp, error) {
	tx := l.r.Begin()

	iss := issuers.Issuer{
		Balance: int(req.Balance),
	}
	if err := tx.Issuers.NewIssuer(&iss); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tx.Commit()
	return &api.NewIssuerResp{
		IssuerID: iss.ID,
	}, nil
}

func (l *Ledger) SellInvoice(ctx context.Context, req *api.SellInvoiceReq) (*api.SellInvoiceResp, error) {
	tx := l.r.Begin()

	if req.InvoiceValue <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invoice value should be > 0")
	}

	issuer := issuers.Issuer{ID: req.IssuerID}
	if err := tx.Issuers.GetIssuer(&issuer); err != nil {
		return nil, status.Errorf(codes.NotFound, "No issuer with such id %v %v", req.IssuerID, err.Error())
	}

	in := invoices.Invoice{}
	in.Value = req.InvoiceValue
	in.IssuerId = req.IssuerID
	in.OwnerID = req.IssuerID
	in.Status = invoices.Available
	in.PutForSale = time.Now()

	l.log.Infow("Putting invoice for sale",
		"Issuer", issuer,
		"Invoice", in,
	)
	if err := tx.Invoices.SaveInvoice(&in); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err := tx.Commit()
	return &api.SellInvoiceResp{InvoiceID: int32(in.ID)}, err
}

func (l *Ledger) GetInvoice(ctx context.Context, req *api.GetInvoiceReq) (*api.Invoice, error) {
	invoice, err := l.r.Begin().Invoices.GetInvoice(req.InvoiceID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &api.Invoice{
		ID:       int32(invoice.ID),
		Value:    int32(invoice.Value),
		IssuerID: invoice.IssuerId,
		OwnerID:  invoice.OwnerID,
		Status:   string(invoice.Status),
	}, nil
}

func (l *Ledger) ListInvoices(ctx context.Context, req *api.ListInvoicesReq) (*api.ListInvoicesResp, error) {
	invoices := l.r.Begin().Invoices.ListInvoices()
	invoicesList := []*api.Invoice{}
	for _, invoice := range invoices {
		invoicesList = append(invoicesList, &api.Invoice{
			ID:       int32(invoice.ID),
			Value:    int32(invoice.Value),
			IssuerID: invoice.IssuerId,
			OwnerID:  invoice.OwnerID,
			Status:   string(invoice.Status),
		})
	}
	return &api.ListInvoicesResp{
		InvoicesList: invoicesList,
	}, nil
}

func (l *Ledger) NewInvestor(ctx context.Context, req *api.NewInvestorReq) (*api.NewInvestorResp, error) {
	tx := l.r.Begin()

	investor := investors.Investor{
		Balance: req.Balance,
	}
	if err := tx.Investors.NewInvestor(&investor); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tx.Commit()
	return &api.NewInvestorResp{
		InvestorId: investor.ID,
	}, nil
}

func (l *Ledger) PlaceBid(ctx context.Context, req *api.PlaceBidReq) (*api.PlaceBidResp, error) {
	tx := l.r.Begin()

	investor, err := tx.Investors.GetInvestor(req.InvestorID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	if investor.Balance-req.Amount < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"can not bid %d, with %d balance. Not enough funds", req.Amount, investor.Balance)
	}
	invoice, err := tx.Invoices.GetInvoice(req.InvoiceID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	if invoice.Status != invoices.Available {
		return nil, status.Error(codes.PermissionDenied, "Invoice is not available for sale")
	}
	err = tx.Investors.ReserveBalance(&investor, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Matching Algorithm
	// After the matching algorithm executes, it releases to the investor the part of the reserved
	// balance that wasn’t used, then investor-3 available balance increases by €200.

	transactionID, err := tx.Transactions.CreateTransaction(transactions.Transaction{
		Amount:     invoice.Value,
		Status:     transactions.Pending,
		InvoiceID:  req.InvoiceID,
		IssuerID:   invoice.IssuerId,
		InvestorID: req.InvestorID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if invoice.Value > req.Amount {
		tx.Investors.ReleaseBalance(&investor, req.Amount)
		tx.Transactions.UpdateTransaction(transactions.Transaction{
			ID:     transactionID,
			Status: transactions.Rejected,
		})
		tx.Commit()
		return nil, status.Errorf(codes.InvalidArgument, "invoice value is %v > bid amount %v", invoice.Value, req.Amount)
	}

	// If bid amount is not the exact value of an invoice - release surplus to the available balance
	if invoice.Value < req.Amount {
		tx.Investors.ReleaseBalance(&investor, req.Amount-invoice.Value)
	}
	invoice.Status = invoices.Financed
	if err := tx.Invoices.UpdateInvoice(invoice); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err = tx.Commit(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.PlaceBidResp{
		Msg:           fmt.Sprintf("Financed invoice %v successfully. Please Approve or Reverse transaction", req.InvoiceID),
		TransactionID: int32(transactionID),
	}, nil
}

func (l *Ledger) ApproveFinancing(ctx context.Context, req *api.ApproveReq) (*api.ApproveResp, error) {
	tx := l.r.Begin()

	l.log.Infof("Trying to approve financing on transaction %v", req.TransactionID)
	transaction, err := tx.Transactions.GetTransaction(req.TransactionID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	l.log.Infof("Transaction to be approved: %+v", transaction)

	if transaction.Status != transactions.Pending {
		return nil, status.Errorf(codes.InvalidArgument, "Transaction %v is already %v", transaction.ID, transaction.Status)
	}
	investor, err := tx.Investors.GetInvestor(transaction.InvestorID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	l.log.Infof("Investor:%+v", investor)
	// 1. Update Invoice
	// 2. Transfer money from investor to issuer
	invoice, err := tx.Invoices.GetInvoice(transaction.InvoiceID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find invoice %v", err)
	}
	l.log.Infof("Invoice to be Approved:%+v", invoice)
	invoice.Status = invoices.Financed
	invoice.OwnerID = investor.ID
	invoice.Financed = time.Now()

	transaction.Status = transactions.Approved
	tx.Transactions.UpdateTransaction(transaction)
	// Transfer Funds from an Investor to the Issuer
	tx.Invoices.UpdateInvoice(invoice)
	tx.Investors.ReduceReservedBalance(&investor, transaction.Amount)
	tx.Issuers.IncreaseBalance(transaction.IssuerID, transaction.Amount)
	tx.Commit()
	return &api.ApproveResp{
		Msg: fmt.Sprintf("Approved financing on invoice %v", invoice.ID),
	}, nil
}

func (l *Ledger) ReverseFinancing(ctx context.Context, req *api.ReverseReq) (*api.ReverseResp, error) {
	tx := l.r.Begin()
	transaction, err := tx.Transactions.GetTransaction(req.TransactionID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if transaction.Status != transactions.Pending {
		return nil, status.Errorf(codes.InvalidArgument, "Transaction %v is already %v", transaction.ID, transaction.Status)
	}
	investor, err := tx.Investors.GetInvestor(transaction.InvestorID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Release investor balance, and revert transaction
	tx.Investors.ReleaseBalance(&investor, transaction.Amount)
	transaction.Status = transactions.Reversed
	tx.Transactions.UpdateTransaction(transaction)

	// Mark invoice as Available for sale, and return it to issuer
	invoice, err := tx.Invoices.GetInvoice(transaction.InvoiceID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find invoice %v", transaction.InvoiceID)
	}
	invoice.Status = invoices.Available
	invoice.OwnerID = transaction.IssuerID
	
	tx.Invoices.UpdateInvoice(invoice)

	if err := tx.Commit(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.ReverseResp{
		Msg: fmt.Sprintf("Reversed financing on invoice %v", transaction.InvoiceID),
	}, nil
}

func (l *Ledger) GetInvestor(ctx context.Context, req *api.GetInvestorReq) (*api.Investor, error) {
	tx := l.r.Begin()
	defer tx.Commit()

	investor, err := tx.Investors.GetInvestor(req.InvestorID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not fund investor with id %v", req.InvestorID)
	}
	invoices := tx.Invoices.ListInvestorInvoices(investor.ID)
	transactions := tx.Transactions.ListInvestorTransactions(investor.ID)
	return &api.Investor{
		ID:              investor.ID,
		Balance:         int32(investor.Balance),
		ReservedBalance: int32(investor.ReservedBalance),
		Invoices:        convertInvoices(invoices),
		Transactions:    convertTransactions(transactions),
	}, nil

}

func (l *Ledger) ListInvestors(ctx context.Context, req *api.ListInvestorsReq) (*api.ListInvestorsResp, error) {
	tx := l.r.Begin()
	defer tx.Commit()

	investorList := tx.Investors.ListInvestors()
	resp := []*api.Investor{}
	for _, investor := range investorList {
		invoices := tx.Invoices.ListInvestorInvoices(investor.ID)
		transactions := tx.Transactions.ListInvestorTransactions(investor.ID)
		resp = append(resp, &api.Investor{
			ID:              investor.ID,
			Balance:         int32(investor.Balance),
			ReservedBalance: int32(investor.ReservedBalance),
			Invoices:        convertInvoices(invoices),
			Transactions:    convertTransactions(transactions),
		})
	}
	return &api.ListInvestorsResp{
		Investors: resp,
	}, nil
}

// Convert database response to the api response
func convertTransactions(transactionsDB []transactions.Transaction) []*api.Transaction {
	transactionsList := []*api.Transaction{}
	for _, transaction := range transactionsDB {
		transactionsList = append(transactionsList, &api.Transaction{
			ID:         int32(transaction.ID),
			Amount:     int32(transaction.Amount),
			Status:     string(transaction.Status),
			InvoiceID:  int32(transaction.InvoiceID),
			IssuerID:   transaction.IssuerID,
			InvestorID: int32(transaction.InvestorID),
			CreatedAt:  transaction.CreatedAt.String(),
			UpdatedAt:  transaction.UpdatedAt.String(),
		})
	}
	return transactionsList
}

// Convert database response to the api response
func convertInvoices(invoicesDB []invoices.Invoice) []*api.Invoice {
	invoicesList := []*api.Invoice{}
	for _, invoice := range invoicesDB {
		invoicesList = append(invoicesList, &api.Invoice{
			ID:       invoice.ID,
			Value:    invoice.Value,
			IssuerID: invoice.IssuerId,
			OwnerID:  invoice.OwnerID,
			Status:   string(invoice.Status),
		})
	}
	return invoicesList
}
