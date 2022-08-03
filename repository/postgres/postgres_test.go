package postgres

import (
	"testing"

	"github.com/Monrevil/simplified-market-ledger/transactions"
)

func TestPostgresDatabase(t *testing.T) {
	r := NewPostgresRepository()
	tx := r.Begin()
	
	tr := transactions.Transaction{
		Amount:     100,
		Status:     0,
		InvoiceID:  0,
		IssuerID:   1,
		InvestorID: 10,
	}
	id, err := tx.Transactions.CreateTransaction(tr)
	if err != nil {
		t.Fatal(err)
	}
	tr2, err := tx.Transactions.GetTransaction(uint(id))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tr2)
}
