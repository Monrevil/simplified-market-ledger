package main

import (
	"fmt"
	"time"

	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/repository/postgres"
	"github.com/Monrevil/simplified-market-ledger/transactions"
)

type Ledger struct {
	r postgres.PostgresRepository
}

func (l *Ledger) SellInvoice(iss issuers.Issuer, in invoices.Invoice) {
	tx := l.r.Begin()

	in.Status = invoices.Available
	in.PutForSale = time.Now()
	tx.Invoices.SaveInvoice(in)
}

func (l *Ledger) PlaceBid(investorID uint, invoiceID int, amount int) error {
	tx := l.r.Begin()

	investor, err := tx.Investors.GetInvestor(investorID)
	if err != nil {
		return err
	}
	if investor.Balance-amount < 0 {
		return fmt.Errorf("can not bid %d, with %d balance. Not enough funds", amount, investor.Balance)
	}

	err = tx.Investors.ReserveBalance(investorID, amount)
	if err != nil {
		return err
	}

	// Matching Algorithm
	// After the matching algorithm executes, it releases to the investor the part of the reserved
	// balance that wasn’t used, then investor-3 available balance increases by €200.
	invoice, err := tx.Invoices.GetInvoice(invoiceID)
	if err != nil {
		return err
	}
	transactionID, err := tx.Transactions.CreateTransaction(transactions.Transaction{
		Amount:     amount,
		Status:     transactions.Pending,
		InvoiceID:  uint(invoiceID),
		IssuerID:   invoice.IssuerId,
		InvestorID: investorID,
	})
	if err != nil {
		return err
	}

	if invoice.Value > amount || invoice.Status != invoices.Available {
		// TODO: Should be transactional
		tx.Investors.ReleaseBalance(investorID, amount)
		tx.Transactions.UpdateTransaction(transactions.Transaction{
			ID:       uint(transactionID),
			Status:   transactions.Rejected,
			UpdateAt: time.Now(),
		})
		return fmt.Errorf("invoice value is %v > bid amount %v", invoice.Value, amount)
	}

	// If bid amount is not the exact value of an invoice - release surplus to the available balance
	if invoice.Value != amount {
		tx.Investors.ReleaseBalance(investorID, amount-invoice.Value)
	}

	return nil
}

func (l *Ledger) Approve(transactionID uint) error {
	tx := l.r.Begin()

	transaction, err := tx.Transactions.GetTransaction(transactionID)
	if err != nil {
		return err
	}
	investor, err := tx.Investors.GetInvestor(transaction.InvestorID)
	if err != nil {
		return err
	}
	// 1. Update Invoice
	// 2. Transfer money from investor to transaction
	invoice := invoices.Invoice{
		ID:       transaction.InvoiceID,
		Status:   invoices.Financed,
		OwnerID:  investor.ID,
		Owner:    fmt.Sprintf("%s %s", investor.FirstName, investor.LastName),
		Financed: time.Now(),
	}
	transaction.Status = transactions.Approved

	tx.Transactions.UpdateTransaction(transaction)
	// TODO Should be Transactional:
	tx.Invoices.UpdateInvoice(invoice)
	tx.Investors.ChangeReservedBalance(transaction.InvestorID, transaction.Amount)
	tx.Issuers.ChangeBalance(transaction.IssuerID, transaction.Amount)

	return nil
}

func (l *Ledger) Reverse(transactionID uint) error {
	tx := l.r.Begin()
	transaction, err := tx.Transactions.GetTransaction(transactionID)
	if err != nil {
		return err
	}

	// TODO Should be Transactional:
	tx.Investors.ReleaseBalance(transaction.InvestorID, transaction.Amount)
	transaction.Status = transactions.Reversed
	tx.Transactions.UpdateTransaction(transaction)
	tx.Invoices.UpdateInvoice(invoices.Invoice{ID: transaction.InvoiceID, Status: invoices.Available})

	return nil
}
