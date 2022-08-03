package postgres

import (
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"gorm.io/gorm"
)

type PostgresInvoicesRepository struct {
	db *gorm.DB
}

func (p *PostgresInvoicesRepository) SaveInvoice(invoice *invoices.Invoice) error {
	return p.db.Create(&invoice).Error
}
func (p *PostgresInvoicesRepository) GetInvoice(invoiceID uint) (invoices.Invoice, error) {
	invoice := invoices.Invoice{}
	err := p.db.First(&invoice, invoiceID).Error
	return invoice, err
}
func (p *PostgresInvoicesRepository) UpdateInvoice(invoice invoices.Invoice) error {
	return p.db.Save(&invoice).Error
}

func (p *PostgresInvoicesRepository) ListInvoices() []invoices.Invoice {
	invoices := []invoices.Invoice{}
	p.db.Find(&invoices)
	return invoices
}

func (p *PostgresInvoicesRepository) ListInvestorInvoices(investorID int32) []invoices.Invoice {
	invoicesList := []invoices.Invoice{}
	p.db.Where(invoices.Invoice{OwnerID: investorID, Status: invoices.Financed}).Find(&invoicesList)
	return invoicesList
}
