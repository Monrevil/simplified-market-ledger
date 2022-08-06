package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Monrevil/simplified-market-ledger/api"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var addr = "localhost:50051"

func TestMain(m *testing.M) {
	// go func() {
	// 	Serve(addr)
	// }()

	var err error
	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	code := m.Run()
	conn.Close()

	os.Exit(code)
}

// Test should be run after docker compose up. With ledger
// TODO: in order to run without using docker compose
// 		 implement using https://github.com/ory/dockertest#using-dockertest
func TestLedger(t *testing.T) {
	c := api.NewLedgerClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()

	issuer, err := c.NewIssuer(ctx, &api.NewIssuerReq{
		Balance: 0,
	})
	require.NoError(t, err)

	invoiceValue := int32(100)
	soldInvoice, err := c.SellInvoice(ctx, &api.SellInvoiceReq{
		IssuerID:     issuer.IssuerID,
		InvoiceValue: invoiceValue,
	})
	require.NoError(t, err)
	invoiceDB, err := c.GetInvoice(ctx, &api.GetInvoiceReq{
		InvoiceID: soldInvoice.InvoiceID,
	})
	require.NoError(t, err)
	require.Equal(t, invoiceValue, invoiceDB.Value, "Invoice value should be recorded in DB")
	require.Equal(t, string(invoices.Available), invoiceDB.Status, "Invoice should be available for sale")
	require.Equal(t, issuer.IssuerID, invoiceDB.OwnerID, "Invoice should belong to the issuer")

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

	investorDB, err := c.GetInvestor(ctx, &api.GetInvestorReq{
		InvestorID: investor.InvestorId,
	})
	require.NoError(t, err)
	require.Equal(t, int32(900), investorDB.Balance, "Should reserve balance after placed bid")
	require.Equal(t, int32(100), investorDB.ReservedBalance, "Should reserve balance after placed bid")

	financed, err := c.ApproveFinancing(ctx, &api.ApproveReq{
		TransactionID: bid.TransactionID,
	})
	require.NoError(t, err)

	boughtInvoice, err := c.GetInvoice(ctx, &api.GetInvoiceReq{
		InvoiceID: soldInvoice.InvoiceID,
	})
	require.NoError(t, err)
	require.Equal(t, investor.InvestorId, boughtInvoice.OwnerID, "Investor should get the invoice")

	t.Log(financed.Msg)
}
