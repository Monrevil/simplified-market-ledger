package invoices

import "time"

type Invoice struct {
	ID         int32
	Value      int32
	Status     InvoiceStatus
	IssuerId   int32
	OwnerID    int32
	PutForSale time.Time
	Financed   time.Time
}

type InvoiceStatus string

const (
	Available InvoiceStatus = "Available"
	Financed  InvoiceStatus = "Financed"
	// Money should be reserved until operation is approved
	Committed InvoiceStatus = "Committed"
)
