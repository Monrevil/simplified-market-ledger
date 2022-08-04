package main

import (
	"testing"
	"time"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
)

func TestE2E(t *testing.T) {
	ledger := Ledger{
		r: postgres.NewPostgresRepository(),
	}

	iss := GetTestIssuer()
	investor := GetTestInvestor()
	tx := ledger.r.Begin()
	tx.Issuers.NewIssuer(iss)
	tx.Investors.NewInvestor(investor)
	tx.Rollback()
}

func GetTestIssuer() *issuers.Issuer {
	return &issuers.Issuer{
		FirstName: "Issuer-1",
		LastName:  "Pangolin",
		Balance:   100,
	}
}

func GetTestInvestor() *investors.Investor {
	return &investors.Investor{
		FirstName:       "Investor",
		LastName:        "Albacore",
		Balance:         1000,
		ReservedBalance: 200,
	}
}

func GetTestInvoice() *invoices.Invoice {
	return &invoices.Invoice{
		Value:      100,
		Status:     invoices.Available,
		IssuerId:   1,
		OwnerID:    1,
		PutForSale: time.Now(),
	}
}
