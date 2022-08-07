package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Monrevil/simplified-market-ledger/ledger/investors"
	"github.com/Monrevil/simplified-market-ledger/ledger/invoices"
	"github.com/Monrevil/simplified-market-ledger/ledger/issuers"
	"github.com/Monrevil/simplified-market-ledger/ledger/transactions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	if port == "" {
		log.Println("POSTGRES_PORT was not set in env, using 5432")
		port = "5432"
	}
	dsn := fmt.Sprintf("host=%s user=test password=test dbname=postgres port=%s sslmode=disable TimeZone=Europe/Kiev", host, port)
	log.Printf("Trying to connect to %v", dsn)

	var db *gorm.DB
	var err error
	err = retry(4, 2*time.Second, func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		return err
	})
	if err != nil {
		panic(err)
	}
	// Migrate the schema
	db.AutoMigrate(&investors.Investor{}, &invoices.Invoice{}, &issuers.Issuer{}, &transactions.Transaction{})

	return &PostgresRepository{
		db: db,
	}
}

func retry(attempts int, interval time.Duration, f func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Printf("retrying after error: %v in %v seconds", err, interval.Seconds())
			time.Sleep(interval)
			interval *= 2
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %v", attempts, err)
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
