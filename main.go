package main

import "log"

func main() {

}

func Run(ledger *Ledger, issuers []Issuer, investors []Investor) {
	for _, iss := range issuers {
		for _, invoice := range iss.Invoices {
			iss.SellInvoice(ledger, invoice)
		}
	}
	for invoiceID := range ledger.Invoices {
		for _, investor := range investors {
			log.Printf("%v balance before bid: %v", investor.FirstName, investor.Balance)
			log.Printf("%v invoices before bid: %v", investor.FirstName, investor.Invoices)
			log.Printf("")
			investor.PlaceBid(ledger, invoiceID, investor.Balance)
			log.Printf("%v balance after bid: %v", investor.FirstName, investor.Balance)
			log.Printf("%v invoices after bid: %v", investor.FirstName, investor.Invoices)
		}
	}
}
