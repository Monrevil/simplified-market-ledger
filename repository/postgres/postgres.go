package postgres

import (
	"fmt"
	"log"
	"os"

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
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		log.Println("POSTGRES_HOST was not set in env, using localhost")
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if host == "" {
		log.Println("POSTGRES_PORT was not set in env, using 5432")
		port = "5432"
	}
	dsn := fmt.Sprintf("host=%s user=test password=test dbname=postgres port=%s sslmode=disable TimeZone=Asia/Shanghai", host, port)
	log.Printf("Trying to connect to %v", dsn)
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
			log.Println("Panicked! performing a db rollback")
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
func (p *PostgresTransaction) Commit() error {
	return p.db.Commit().Error
}

// Investors:
type PostgresInvestorsRepository struct {
	db *gorm.DB
}

func (p *PostgresInvestorsRepository) NewInvestor(investor *investors.Investor) error {
	return p.db.Create(investor).Error
}

func (p *PostgresInvestorsRepository) GetInvestor(investorID int32) (investors.Investor, error) {
	investor := investors.Investor{ID: investorID}
	err := p.db.First(&investor, investorID).Error
	return investor, err
}

// ReserveBalance expects ID to be set for the investor
func (p *PostgresInvestorsRepository) ReserveBalance(investor *investors.Investor, amount int) error {
	if err := p.db.First(investor).Error; err != nil {
		return err
	}

	if investor.Balance < amount {
		return fmt.Errorf("not enough funds. Tired to reserve %v funds. With %v active balance", amount, investor.Balance)
	}

	err := p.db.Model(investor).Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		return err
	}
	err = p.db.Model(investor).Update("reserved_balance", gorm.Expr("reserved_balance + ?", amount)).Error
	if err != nil {
		return err
	}

	if err := p.db.First(investor).Error; err != nil {
		return err
	}

	return nil
}
func (p *PostgresInvestorsRepository) ReleaseBalance(investor *investors.Investor, amount int) error {
	if err := p.db.First(investor).Error; err != nil {
		return err
	}
	if investor.ReservedBalance < amount {
		return fmt.Errorf("not enough funds. Tired to release %v funds. With %v reserved balance", amount, investor.ReservedBalance)
	}

	err := p.db.Model(investor).Update("reserved_balance", gorm.Expr("reserved_balance - ?", amount)).Error
	if err != nil {
		return err
	}
	err = p.db.Model(investor).Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return err
	}
	if err := p.db.First(investor).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresInvestorsRepository) ReduceReservedBalance(investor *investors.Investor, amount int) error {
	if investor.ID == 0 {
		return fmt.Errorf("please provide non-0 investor id")
	}
	investor.ReservedBalance -= amount
	err := p.db.Model(investor).Update("reserved_balance", gorm.Expr("reserved_balance - ?", amount)).Error

	return err
}

func (p *PostgresInvestorsRepository) ListInvestors() []investors.Investor {
	var investors []investors.Investor
	p.db.Find(&investors)
	return investors
}
