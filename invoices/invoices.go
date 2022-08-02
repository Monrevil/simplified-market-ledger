package invoices

import "time"

type Repository interface {
	SaveInvoice(Invoice) error
	GetInvoice(int) (Invoice, error)
	UpdateInvoice(Invoice) error
}

type Invoice struct {
	ID         uint
	Value      int
	Status     InvoiceStatus
	Issuer     string
	IssuerId   uint
	OwnerID    uint
	Owner      string
	PutForSale time.Time
	Financed   time.Time
}

type InvoiceStatus uint8

const (
	Available InvoiceStatus = 0
	Financed  InvoiceStatus = 1
	Reversed  InvoiceStatus = 2
	// Money should be reserved until operation is approved
	Committed InvoiceStatus = 3
)
