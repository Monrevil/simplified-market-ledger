package postgres

import (
	"github.com/Monrevil/simplified-market-ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/transactions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository() *PostgresRepository {
	dsn := "host=localhost user=test password=test dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&investors.Investor{}, &invoices.Invoice{}, &issuers.Issuer{}, &transactions.Transaction{})

	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Begin() *PostgresTransaction {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	return &PostgresTransaction{
		db: tx,
		Investors: PostgresInvestorsRepository{
			db: tx,
		},
		Issuers: PostgresIssuersRepository{
			db: tx,
		},
		Invoices: PostgresInvoicesRepository{
			db: tx,
		},
		Transactions: PostgresTransactionsRepository{
			db: tx,
		},
	}
}

type PostgresTransaction struct {
	db           *gorm.DB
	Investors    PostgresInvestorsRepository
	Issuers      PostgresIssuersRepository
	Invoices     PostgresInvoicesRepository
	Transactions PostgresTransactionsRepository
}

// Rollback a transaction
func (p *PostgresTransaction) Rollback() {
	p.db.Rollback()
}

// Commit a transaction
func (p *PostgresTransaction) Commit() {
	p.db.Commit()
}

// Investors:
type PostgresInvestorsRepository struct {
	db *gorm.DB
}

func (p *PostgresInvestorsRepository) GetInvestor(investorID uint) (investors.Investor, error) {
	investor := investors.Investor{}
	err := p.db.First(&investor, investorID).Error
	return investor, err
}

func (p *PostgresInvestorsRepository) ReserveBalance(investorID uint, amount int) error {
	investor := &investors.Investor{ID: investorID}
	err := p.db.Model(investor).Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		return err
	}
	p.db.Model(investor).Update("reserved_balance", gorm.Expr("balance + ?", amount))
	if err != nil {
		return err
	}
	return nil
}
func (p *PostgresInvestorsRepository) ReleaseBalance(investorID uint, amount int) error {
	investor := &investors.Investor{ID: investorID}
	err := p.db.Model(investor).Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return err
	}
	p.db.Model(investor).Update("reserved_balance", gorm.Expr("balance - ?", amount))
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresInvestorsRepository) ChangeReservedBalance(investorID uint, amount int) error {
	investor := &investors.Investor{ID: investorID}
	err := p.db.Model(investor).Update("balance", gorm.Expr("balance + ?", amount)).Error

	return err
}

func (p *PostgresInvestorsRepository) ListInvestors() []investors.Investor {
	var investors []investors.Investor
	p.db.Find(&investors)
	return investors
}

// Invoices:

type PostgresInvoicesRepository struct {
	db *gorm.DB
}

func (p *PostgresInvoicesRepository) SaveInvoice(invoice invoices.Invoice) error {
	return p.db.Save(&invoice).Error
}
func (p *PostgresInvoicesRepository) GetInvoice(invoiceID int) (invoices.Invoice, error) {
	invoice := invoices.Invoice{}
	err := p.db.First(&invoice, invoiceID).Error
	return invoice, err
}
func (p *PostgresInvoicesRepository) UpdateInvoice(invoice invoices.Invoice) error {
	return p.db.Save(&invoice).Error
}

// Issuers:

type PostgresIssuersRepository struct {
	db *gorm.DB
}

func (p *PostgresIssuersRepository) ChangeBalance(id uint, amount int) error {
	issuer := &issuers.Issuer{ID: id}
	return p.db.Model(issuer).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// Transactions:

type PostgresTransactionsRepository struct {
	db *gorm.DB
}

func (p *PostgresTransactionsRepository) GetTransaction(transactionID uint) (transactions.Transaction, error) {
	tr := transactions.Transaction{ID: transactionID}
	err := p.db.First(&tr).Error
	return tr, err
}

func (p *PostgresTransactionsRepository) CreateTransaction(tr transactions.Transaction) (int, error) {
	result := p.db.Create(&tr)
	return int(tr.ID), result.Error
}
func (p *PostgresTransactionsRepository) UpdateTransaction(tr transactions.Transaction) error {
	return p.db.Save(&tr).Error
}
