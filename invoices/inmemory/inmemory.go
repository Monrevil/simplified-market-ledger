package inmemory

import (
	"fmt"

	"github.com/Monrevil/simplified-market-ledger/invoices"
)

//  This package is used for storing data in memory, data is not persistent.
//  Main purpose of this package is testing
type MemoryInvoiceRepository struct {
	invoices map[uint]invoices.Invoice
}

func NewMemoryInvoiceRepository() *MemoryInvoiceRepository {
	return &MemoryInvoiceRepository{}
}

func (m *MemoryInvoiceRepository) Save(in invoices.Invoice) error {
	newID := uint(len(m.invoices))
	in.ID = newID
	m.invoices[newID] = in
	return nil
}

func (m *MemoryInvoiceRepository) Get(invoiceID int) (invoices.Invoice, error) {
	invoice, ok := m.invoices[uint(invoiceID)]
	if !ok {
		return invoices.Invoice{}, fmt.Errorf("Not found")
	}
	return invoice, nil
}
