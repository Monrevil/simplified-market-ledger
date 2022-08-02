package main

import "fmt"

type Investor struct {
	FirstName string
	LastName  string
	Balance   int
	Invoices  []Invoice
}

func (i *Investor) PlaceBid(l *Ledger, invoiceID uint, amount int) error {
	if i.Balance-amount < 0 {
		return fmt.Errorf("can not bid %d, with %d balance. Not enough funds", amount, i.Balance)
	}
	i.Balance -= amount
	invoice, err := l.TakeBid(i, invoiceID, amount)
	if err != nil {
		i.Balance += amount
		return err
	}
	i.Invoices = append(i.Invoices, invoice)
	return nil
}
