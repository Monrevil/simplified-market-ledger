package main

type Issuer struct {
	FirstName string
	LastName  string
	Balance   int
	Invoices  []Invoice
}

func (iss *Issuer) SellInvoice(l *Ledger, in Invoice) {
	
	l.SellInvoice(in)
}