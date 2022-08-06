package postgres_test

import (
	"testing"
	"time"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
	"github.com/Monrevil/simplified-market-ledger/transactions"
	"github.com/stretchr/testify/require"
)

func TestPostgresDatabase(t *testing.T) {
	r := postgres.NewPostgresRepository()

	tr := transactions.Transaction{
		Amount:     100,
		Status:     transactions.Approved,
		InvoiceID:  0,
		IssuerID:   1,
		InvestorID: 10,
	}

	t.Run("Transactions repository", func(t *testing.T) {
		tx := r.Begin()

		id, err := tx.Transactions.CreateTransaction(tr)
		require.Nil(t, err)

		tr2, err := tx.Transactions.GetTransaction(uint(id))
		require.Nil(t, err)

		tr2.Status = transactions.Rejected
		err = tx.Transactions.UpdateTransaction(tr2)
		require.Nil(t, err)

		tr2, err = tx.Transactions.GetTransaction(uint(id))
		require.Nil(t, err)
		require.Equal(t, transactions.Rejected, tr2.Status)

		tx.Rollback()
	})

	t.Run("Invoices repository", func(t *testing.T) {
		tx := r.Begin()
		in := invoices.Invoice{
			Value:      0,
			Status:     invoices.Available,
			IssuerId:   1,
			OwnerID:    1,
			PutForSale: time.Now(),
		}
		err := tx.Invoices.SaveInvoice(&in)
		require.Nil(t, err)

		in.Status = invoices.Financed
		err = tx.Invoices.UpdateInvoice(in)
		require.Nil(t, err)

		in, err = tx.Invoices.GetInvoice(in.ID)
		require.Nil(t, err)
		require.Equal(t, invoices.Financed, in.Status)

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
		require.Equal(t, 200, iss.Balance)

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
	err := tx.Investors.NewInvestor(i)
	require.NoError(t, err)
	require.NotEqual(t, 0, i.ID, "Should set investor ID")

	err = tx.Investors.ReleaseBalance(i, 100)
	require.NoError(t, err)
	require.Equal(t, 1100, i.Balance, "Releasing balance should change active balance")
	require.Equal(t, 100, i.ReservedBalance, "Releasing balance should change reserved balance")

	err = tx.Investors.ReserveBalance(i, 600)
	require.NoError(t, err)
	require.Equal(t, 500, i.Balance, "Reserving balance should change active balance")
	require.Equal(t, 700, i.ReservedBalance, "Reserving balance should change reserved balance")

	err = tx.Investors.ReserveBalance(i, 600)
	require.Error(t, err, "Should return error if active balance < amount to be reserved")

	err = tx.Investors.ReduceReservedBalance(i, 600)
	require.Nil(t, err)

	*i, err = tx.Investors.GetInvestor(i.ID)
	require.Nil(t, err)
	require.Equal(t, 100, i.ReservedBalance)

	err = tx.Investors.ReduceReservedBalance(&investors.Investor{ID: 0}, 600)
	require.Error(t, err, "Should return error for non-existing investor")

	err = tx.Investors.ReserveBalance(&investors.Investor{ID: 0}, 600)
	require.Error(t, err, "Should return error for non-existing investor")

	*i, err = tx.Investors.GetInvestor(i.ID)
	require.Nil(t, err)
	require.Equal(t, 100, i.ReservedBalance)

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
		FirstName:       "Investor-1",
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
