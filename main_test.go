package main

import (
	"fmt"
	"testing"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/issuers"
)

func TestE2E(t *testing.T) {
	issuers := []issuers.Issuer{
		{
			FirstName: "Investor1",
			LastName:  "Johnson",
			Balance:   0,
		},
	}
	investors := []investors.Investor{
		{
			FirstName: "Investor1",
			LastName:  "Johnson",
			Balance:   10,
		},
		{
			FirstName: "Investor2",
			LastName:  "Johnson",
			Balance:   110,
		},
		{
			FirstName: "Investor3",
			LastName:  "Johnson",
			Balance:   200,
		},
	}

	ledger := Ledger{
		invoicesRepository:    nil,
		issuersRepository:     nil,
		investorsRepository:   nil,
		transactionRepository: nil,
	}

	fmt.Println(&ledger, issuers, investors)
}
