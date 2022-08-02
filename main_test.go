package main

import (
	"testing"
	"time"
)

func TestE2E(t *testing.T) {
	issuers := []Issuer{
		{
			FirstName: "Investor1",
			LastName:  "Johnson",
			Balance:   0,
			Invoices: []Invoice{
				{
					Id:         0,
					Value:      100,
					Status:     0,
					OwnerID:    "",
					Owner:      "Bob1",
					PutForSale: time.Time{},
					Financed:   time.Time{},
				},
			},
		},
	}
	investors := []Investor{
		{
			FirstName: "Investor1",
			LastName:  "Johnson",
			Balance:   10,
			Invoices:  []Invoice{},
		},
		{
			FirstName: "Investor2",
			LastName:  "Johnson",
			Balance:   110,
			Invoices:  []Invoice{},
		},
		{
			FirstName: "Investor3",
			LastName:  "Johnson",
			Balance:   200,
			Invoices:  []Invoice{},
		},

	}

	ledger := Ledger{
		Invoices: map[uint]Invoice{},
	}
	// for _, iss := range issuers {
	// 	for _, invoice := range iss.Invoices {
	// 		iss.SellInvoice(&ledger, invoice)
	// 	}
	// }
	// for invoiceID := range ledger.Invoices {
	// 	for _, investor := range investors {
	// 		t.Logf("%v balance before bid", investor.Balance)
	// 		investor.PlaceBid(&ledger, invoiceID, investor.Balance)
	// 		t.Logf("%v balance after bid", investor.Balance)
	// 	}
	// }

	Run(&ledger, issuers, investors)
}
