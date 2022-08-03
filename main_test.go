package main

import (
	"fmt"
	"testing"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
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
		r: *postgres.NewPostgresRepository(),
	}

	fmt.Println(&ledger, issuers, investors)
}
