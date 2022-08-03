package postgres_test

import (
	"testing"
	"time"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
	"github.com/Monrevil/simplified-market-ledger/transactions"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDatabase(t *testing.T) {
	r := postgres.NewPostgresRepository()

	tr := transactions.Transaction{
		Amount:     100,
		Status:     0,
		InvoiceID:  0,
		IssuerID:   1,
		InvestorID: 10,
	}

	t.Run("Transactions repository", func(t *testing.T) {
		tx := r.Begin()

		id, err := tx.Transactions.CreateTransaction(tr)
		assert.Nil(t, err)
	
		tr2, err := tx.Transactions.GetTransaction(uint(id))
		assert.Nil(t, err)

		tr2.Status = transactions.Rejected
		err = tx.Transactions.UpdateTransaction(tr2)
		assert.Nil(t, err)

		tr2, err = tx.Transactions.GetTransaction(uint(id))
		assert.Nil(t, err)
		assert.Equal(t, transactions.Rejected, tr2.Status)


		tx.Rollback()
	})

	t.Run("Invoices repository", func(t *testing.T) {
		tx := r.Begin()
		in := invoices.Invoice{
			Value:      0,
			Status:     invoices.Available,
			Issuer:     "Issuer-1",
			IssuerId:   1,
			OwnerID:    1,
			Owner:      "Issuer-1",
			PutForSale: time.Now(),
		}
		err := tx.Invoices.SaveInvoice(&in)
		assert.Nil(t, err)

		in.Status = invoices.Financed
		err = tx.Invoices.UpdateInvoice(in)
		assert.Nil(t, err)

		in, err = tx.Invoices.GetInvoice(in.ID)
		assert.Nil(t, err)
		assert.Equal(t, invoices.Financed, in.Status)

		tx.Rollback()
	})

	t.Run("Issuers repository", func(t *testing.T) {
		tx := r.Begin()
		iss := issuers.Issuer{
			FirstName: "Issuer-1",
			LastName:  "Pangolin",
			Balance:   100,
		}

		if err := tx.Issuers.NewIssuer(&iss); err != nil {
			t.Fatal(err)
		}
		if err := tx.Issuers.ChangeBalance(iss.ID, 100); err != nil {
			t.Fatal(err)
		}
		if err := tx.Issuers.GetIssuer(&iss); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, iss.Balance)

		tx.Rollback()
	})

}

func TestInvestorsRepository(t *testing.T) {
	r := postgres.NewPostgresRepository()
	tx := r.Begin()

	i := &investors.Investor{
		FirstName:       "Investor-1",
		LastName:        "Albacore",
		Balance:         1000,
		ReservedBalance: 200,
	}
	if err := tx.Investors.NewInvestor(i); err != nil {
		t.Fatal(err)
	}

	if err := tx.Investors.ReleaseBalance(i, 100); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1100, i.Balance, "Releasing balance should change active balance")
	assert.Equal(t, 100, i.ReservedBalance, "Releasing balance should change reserved balance")

	if err := tx.Investors.ReserveBalance(i, 600); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 500, i.Balance, "Reserving balance should change active balance")
	assert.Equal(t, 700, i.ReservedBalance, "Reserving balance should change reserved balance")
	spew.Dump(i)

	err := tx.Investors.ReserveBalance(i, 600)
	assert.NotNil(t, err, "Should return error if active balance < amount to be reserved")

	err = tx.Investors.ReduceReservedBalance(i, 600)
	assert.Nil(t, err)

	*i, err = tx.Investors.GetInvestor(i.ID)
	assert.Nil(t, err)
	assert.Equal(t, 100, i.ReservedBalance)

	err = tx.Investors.ReduceReservedBalance(&investors.Investor{ID: 0}, 600)
	assert.NotNil(t, err, "Should return error for non-existing investor")

	*i, err = tx.Investors.GetInvestor(i.ID)
	assert.Nil(t, err)
	assert.Equal(t, 100, i.ReservedBalance)

	tx.Rollback()
}
