package invoices

import "time"

type Repository interface {
	SaveInvoice(invoice Invoice) error
	GetInvoice(invoiceID int) (Invoice, error)
	UpdateInvoice(invoice Invoice) error
}

type Invoice struct {
	ID         uint
	Value      int
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
