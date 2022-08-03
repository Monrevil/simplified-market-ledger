package transactions

import "time"

type Repository interface {
	GetTransaction(transactionID uint) (Transaction, error)
	CreateTransaction(tr Transaction) (int, error)
	UpdateTransaction(tr Transaction) error
}

type Transaction struct {
	ID        uint
	Amount    int
	Status    Status
	InvoiceID uint
	// IssuerID - issuer here is a party selling invoice. Beneficiary of the transaction.
	IssuerID   int32
	InvestorID uint
	CreatedAt  time.Time
	UpdateAt   time.Time
}

type Status string

const (
	Pending  Status = "Pending"
	Approved Status = "Approved"
	Reversed Status = "Reversed"
	Rejected Status = "Rejected"
)
