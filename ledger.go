package main

import (
	"fmt"
	"time"

	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/transactions"
)

type Ledger struct {
	invoicesRepository    invoices.Repository
	issuersRepository     issuers.Repository
	investorsRepository   investors.Repository
	transactionRepository transactions.Repository
}

func (l *Ledger) SellInvoice(iss issuers.Issuer, in invoices.Invoice) {
	in.Status = invoices.Available
	in.PutForSale = time.Now()
	l.invoicesRepository.SaveInvoice(in)
}

func (l *Ledger) PlaceBid(investorID uint, invoiceID int, amount int) error {
	investor, err := l.investorsRepository.GetInvestor(investorID)
	if err != nil {
		return err
	}
	if investor.Balance-amount < 0 {
		return fmt.Errorf("can not bid %d, with %d balance. Not enough funds", amount, investor.Balance)
	}

	// TODO: should be transactional:
	err = l.investorsRepository.ReserveBalance(investorID, amount)
	if err != nil {
		return err
	}

	// Matching Algorithm
	// After the matching algorithm executes, it releases to the investor the part of the reserved
	// balance that wasn’t used, then investor-3 available balance increases by €200.
	invoice, err := l.invoicesRepository.GetInvoice(invoiceID)
	if err != nil {
		return err
	}
	transactionID, err := l.transactionRepository.CreateTransaction(transactions.Transaction{
		Amount:     amount,
		Status:     transactions.Pending,
		InvoiceID:  uint(invoiceID),
		IssuerID:   invoice.IssuerId,
		InvestorID: investorID,
		CreatedAt:  time.Now(),
		UpdateAt:   time.Now(),
	})
	if err != nil {
		return err
	}

	if invoice.Value > amount || invoice.Status != invoices.Available {
		// TODO: Should be transactional
		l.investorsRepository.ReleaseBalance(investorID, amount)
		l.transactionRepository.UpdateTransaction(transactions.Transaction{
			ID:       uint(transactionID),
			Status:   transactions.Rejected,
			UpdateAt: time.Now(),
		})
		return fmt.Errorf("invoice value is %v > bid amount %v", invoice.Value, amount)
	}

	// If bid amount is not the exact value of an invoice - release surplus to the available balance
	if invoice.Value != amount {
		l.investorsRepository.ReleaseBalance(investorID, amount-invoice.Value)
	}

	return nil
}

func (l *Ledger) Approve(transactionID uint) error {
	transaction, err := l.transactionRepository.GetTransaction(transactionID)
	if err != nil {
		return err
	}
	investor, err := l.investorsRepository.GetInvestor(transaction.InvestorID)
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

	l.transactionRepository.UpdateTransaction(transaction)
	// TODO Should be Transactional:
	l.invoicesRepository.UpdateInvoice(invoice)
	l.investorsRepository.ChangeReservedBalance(transaction.InvestorID, transaction.Amount)
	l.issuersRepository.ChangeBalance(transaction.IssuerID, transaction.Amount)

	return nil
}

func (l *Ledger) Reverse(transactionID uint) error {
	transaction, err := l.transactionRepository.GetTransaction(transactionID)
	if err != nil {
		return err
	}

	// TODO Should be Transactional:
	l.investorsRepository.ReleaseBalance(transaction.InvestorID, transaction.Amount)
	transaction.Status = transactions.Reversed
	l.transactionRepository.UpdateTransaction(transaction)
	l.invoicesRepository.UpdateInvoice(invoices.Invoice{ID: transaction.InvoiceID, Status: invoices.Available})

	return nil
}
