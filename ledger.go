package main

import (
	"fmt"
	"time"
)

type Ledger struct {
	Invoices map[uint]Invoice
}

func (l *Ledger) SellInvoice(in Invoice) {
	in.Status = Stored
	in.Id = uint(len(l.Invoices))
	in.PutForSale = time.Now()
	l.Invoices[in.Id] = in
}

func (l *Ledger) TakeBid(investor *Investor, invoiceID uint, amount int) (Invoice, error) {
	invoice, ok := l.Invoices[invoiceID]
	if !ok {
		return Invoice{}, fmt.Errorf("invoice with ID %d does not exist", invoiceID)
	}
	if amount < invoice.Value {
		return Invoice{}, fmt.Errorf("rejected. Invoice value is > bid amount")
	}
	invoice.Status = Financed
	invoice.Financed = time.Now()
	invoice.Owner = fmt.Sprintf("%s %s", investor.FirstName, investor.LastName)
	delete(l.Invoices, invoiceID)
	
	return invoice, nil
}
