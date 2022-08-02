package main

import "time"

type Invoice struct {
	Id         uint
	Value      int
	Status     InvoiceStatus
	OwnerID    string
	Owner      string
	PutForSale time.Time
	Financed   time.Time
}

type InvoiceStatus uint8

const (
	Stored   InvoiceStatus = 0
	Financed InvoiceStatus = 1
	Reversed InvoiceStatus = 2
	// Money should be reserved until operation is approved
	Committed InvoiceStatus = 3
)
