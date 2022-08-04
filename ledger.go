package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Monrevil/simplified-market-ledger/api"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
	"github.com/Monrevil/simplified-market-ledger/transactions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Ledger struct {
	r *postgres.PostgresRepository
	api.UnimplementedLedgerServer
}

func Serve() {
	addr := fmt.Sprintf(":%d", 50051)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Cannot listen to address %s", addr)
	}
	s := grpc.NewServer()
	ledgerServer := &Ledger{
		r: postgres.NewPostgresRepository(),
	}
	api.RegisterLedgerServer(s, ledgerServer)

	// Perform graceful shutdown if interrupted.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		stop := <-sigCh
		log.Printf("Got %v signal, attempting graceful shutdown", stop)
		s.GracefulStop()
		log.Printf("Shut down gRPC server successfully")
	}()

	log.Printf("gRPC server is listening and serving on %v", addr)
	if err := s.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
	in.IssuerId = req.IssuerID
	in.OwnerID = req.IssuerID
	in.Status = invoices.Available
	in.PutForSale = time.Now()

	if err := tx.Invoices.SaveInvoice(&in); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err := tx.Commit()
	return &api.SellInvoiceResp{InvoiceID: int32(in.ID)}, err
}

func (l *Ledger) GetInvoice(ctx context.Context, req *api.GetInvoiceReq) (*api.Invoice, error) {
	invoice, err := l.r.Begin().Invoices.GetInvoice(uint(req.InvoiceID))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &api.Invoice{
		ID:      int32(invoice.ID),
		Value:   int32(invoice.Value),
		OwnerID: int32(invoice.Value),
		Status:  string(invoice.Status),
	}, nil
}

func (l *Ledger) ListInvoices(ctx context.Context, req *api.ListInvoicesReq) (*api.ListInvoicesResp, error) {
	invoices := l.r.Begin().Invoices.ListInvoices()
	invoicesList := []*api.Invoice{}
	for _, invoice := range invoices {
		invoicesList = append(invoicesList, &api.Invoice{
			ID:      int32(invoice.ID),
			Value:   int32(invoice.Value),
			OwnerID: invoice.OwnerID,
			Status:  string(invoice.Status),
		})
	}
	return &api.ListInvoicesResp{
		InvoicesList: invoicesList,
	}, nil
}

func (l *Ledger) PlaceBid(ctx context.Context, req *api.PlaceBidReq) (*api.PlaceBidResp, error) {
	tx := l.r.Begin()

	investor, err := tx.Investors.GetInvestor(req.InvestorID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	if investor.Balance-int(req.Amount) < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"can not bid %d, with %d balance. Not enough funds", req.Amount, investor.Balance)
	}
	invoice, err := tx.Invoices.GetInvoice(uint(req.InvoiceID))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	if invoice.Status != invoices.Available {
		return nil, status.Error(codes.PermissionDenied, "Invoice is not available for sale")
	}
	err = tx.Investors.ReserveBalance(&investor, int(req.Amount))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Matching Algorithm
	// After the matching algorithm executes, it releases to the investor the part of the reserved
	// balance that wasn’t used, then investor-3 available balance increases by €200.

	transactionID, err := tx.Transactions.CreateTransaction(transactions.Transaction{
		Amount:     int(req.Amount),
		Status:     transactions.Pending,
		InvoiceID:  uint(req.InvoiceID),
		IssuerID:   invoice.IssuerId,
		InvestorID: uint(req.InvestorID),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if invoice.Value > int(req.Amount) {
		// TODO: Should be transactional
		tx.Investors.ReleaseBalance(&investor, int(req.Amount))
		tx.Transactions.UpdateTransaction(transactions.Transaction{
			ID:       uint(transactionID),
			Status:   transactions.Rejected,
			UpdateAt: time.Now(),
		})
		tx.Commit()
		return nil, status.Errorf(codes.InvalidArgument, "invoice value is %v > bid amount %v", invoice.Value, req.Amount)
	}

	// If bid amount is not the exact value of an invoice - release surplus to the available balance
	if invoice.Value != int(req.Amount) {
		tx.Investors.ReleaseBalance(&investor, int(req.Amount)-invoice.Value)
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

	transaction, err := tx.Transactions.GetTransaction(uint(req.TransactionID))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if transaction.Status != transactions.Pending {
		return nil, status.Errorf(codes.InvalidArgument, "Transaction %v is already %v", transaction.ID, transaction.Status)
	}
	investor, err := tx.Investors.GetInvestor(int32(transaction.InvestorID))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	// 1. Update Invoice
	// 2. Transfer money from investor to issuer
	invoice := invoices.Invoice{
		ID:       transaction.InvoiceID,
		Status:   invoices.Financed,
		OwnerID:  investor.ID,
		Financed: time.Now(),
	}
	transaction.Status = transactions.Approved

	tx.Transactions.UpdateTransaction(transaction)
	// TODO Should be Transactional:
	tx.Invoices.UpdateInvoice(invoice)
	tx.Investors.ReduceReservedBalance(&investor, transaction.Amount)
	tx.Issuers.ChangeBalance(transaction.IssuerID, transaction.Amount)
	tx.Commit()
	return &api.ApproveResp{
		Msg: fmt.Sprintf("Approved financing on invoice %v", invoice.ID),
	}, nil
}

func (l *Ledger) ReverseFinancing(ctx context.Context, req *api.ReverseReq) (*api.ReverseResp, error) {
	tx := l.r.Begin()
	transaction, err := tx.Transactions.GetTransaction(uint(req.TransactionID))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	investor, err := tx.Investors.GetInvestor(int32(transaction.InvestorID))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// TODO Should be Transactional:
	tx.Investors.ReleaseBalance(&investor, transaction.Amount)
	transaction.Status = transactions.Reversed
	tx.Transactions.UpdateTransaction(transaction)
	tx.Invoices.UpdateInvoice(invoices.Invoice{ID: transaction.InvoiceID, Status: invoices.Available})

	if err := tx.Commit(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.ReverseResp{
		Msg: fmt.Sprintf("Reversed financing on invoice %v", transaction.InvoiceID),
	}, nil
}

func (l *Ledger) ListInvestors(ctx context.Context, req *api.ListInvestorsReq) (*api.ListInvestorsResp, error) {
	investorList := l.r.Begin().Investors.ListInvestors()
	resp := []*api.Investor{}
	for _, investor := range investorList {
		invoicesDB := l.r.Begin().Invoices.ListInvestorInvoices(investor.ID)
		invoicesList := []*api.Invoice{}
		for _, invoice := range invoicesDB {
			invoicesList = append(invoicesList, &api.Invoice{
				ID:      int32(invoice.ID),
				Value:   int32(invoice.Value),
				OwnerID: invoice.OwnerID,
				Status:  string(invoice.Status),
			})
		}
		transactionsDB := l.r.Begin().Transactions.ListInvestorTransactions(investor.ID)
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
				UpdatedAt:  transaction.UpdateAt.String(),
			})
		}
		resp = append(resp, &api.Investor{
			ID:              investor.ID,
			Balance:         int32(investor.Balance),
			ReservedBalance: int32(investor.ReservedBalance),
			Invoices:        invoicesList,
			Transactions:    transactionsList,
		})
	}
	return &api.ListInvestorsResp{
		Investors: resp,
	}, nil
}
