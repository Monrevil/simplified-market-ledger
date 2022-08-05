package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Monrevil/simplified-market-ledger/api"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Test should be run after docker compose up. With ledger
// TODO: in order to run without using docker compose
// 		 implement using https://github.com/ory/dockertest#using-dockertest
func TestLedger(t *testing.T) {
	// addr := "localhost:5050"
	addr := "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewLedgerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.NewIssuer(ctx, &api.NewIssuerReq{
		Balance: 0,
	})
	require.NoError(t, err)

	invoiceValue := int32(100)
	soldInvoice, err := c.SellInvoice(ctx, &api.SellInvoiceReq{
		IssuerID:     r.IssuerID,
		InvoiceValue: invoiceValue,
	})
	require.NoError(t, err)
	invoiceDB, err := c.GetInvoice(ctx, &api.GetInvoiceReq{
		InvoiceID: soldInvoice.InvoiceID,
	})
	require.NoError(t, err)
	require.Equal(t, invoiceDB.Value, invoiceValue)
	require.Equal(t, invoiceDB.Status, invoices.Available)
	require.Equal(t, invoiceDB.OwnerID, r.IssuerID)

	investor, err := c.NewInvestor(ctx, &api.NewInvestorReq{
		Balance: 1000,
	})
	require.NoError(t, err)

	bid, err := c.PlaceBid(ctx, &api.PlaceBidReq{
		InvestorID: investor.InvestorId,
		InvoiceID:  soldInvoice.InvoiceID,
		Amount:     100,
	})
	require.NoError(t, err)

	financed, err := c.ApproveFinancing(ctx, &api.ApproveReq{
		TransactionID: bid.TransactionID,
	})
	require.NoError(t, err)

	boughtInvoice, err := c.GetInvoice(ctx, &api.GetInvoiceReq{
		InvoiceID: soldInvoice.InvoiceID,
	})
	require.NoError(t, err)
	require.Equal(t, investor.InvestorId, boughtInvoice.OwnerID, "Investor did not get the invoice")

	t.Log(financed.Msg)
}
