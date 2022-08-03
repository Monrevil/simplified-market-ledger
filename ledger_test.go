package main

import (
	"testing"
	"time"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	ledger := Ledger{
		r: postgres.NewPostgresRepository(),
	}

	iss := GetTestIssuer()
	inv := GetTestInvoice()
	investor := GetTestInvestor()
	tx := ledger.r.Begin()
	tx.Issuers.NewIssuer(iss)
	tx.Investors.NewInvestor(investor)
	tx.Commit()

	invoiceID, err := ledger.SellInvoice(*iss, *inv)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)

	tx = ledger.r.Begin()
	*inv, err = tx.Invoices.GetInvoice(invoiceID)
	assert.Nil(t, err)
	assert.Equal(t, invoices.Available, inv.Status)
	spew.Dump(inv)
	tx.Commit()
	
	err = ledger.PlaceBid(investor.ID, invoiceID, 200)
	assert.Nil(t, err)

	investors := ledger.ListInvestors()
	spew.Dump(investors)

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
		Issuer:     "Issuer-1",
		IssuerId:   1,
		OwnerID:    1,
		Owner:      "Issuer-1",
		PutForSale: time.Now(),
	}
}
