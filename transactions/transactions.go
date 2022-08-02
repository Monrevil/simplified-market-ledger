package transactions

import "time"

type Repository interface {
	GetTransaction(transactionID uint) (Transaction, error)
	CreateTransaction(Transaction) (int, error)
	UpdateTransaction(Transaction) error
}

type Transaction struct {
	ID        uint
	Amount    int
	Status    Status
	InvoiceID uint
	// IssuerID - issuer here is a party selling invoice. Beneficiary of the transaction.
	IssuerID   uint
	InvestorID uint
	CreatedAt  time.Time
	UpdateAt   time.Time
}

type Status uint

const (
	Pending  Status = 0
	Approved Status = 1
	Reversed Status = 2
	Rejected Status = 3
)
