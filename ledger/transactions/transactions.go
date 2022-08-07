package transactions

import "time"

type Transaction struct {
	ID        int32
	Amount    int32
	Status    Status
	InvoiceID int32
	// IssuerID - issuer here is a party selling invoice. Beneficiary of the transaction.
	IssuerID   int32
	InvestorID int32
	CreatedAt  time.Time
	UpdatedAt   time.Time
}

type Status string

const (
	Pending  Status = "Pending"
	Approved Status = "Approved"
	Reversed Status = "Reversed"
	Rejected Status = "Rejected"
)
