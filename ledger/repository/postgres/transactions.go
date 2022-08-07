package postgres

import (
	"github.com/Monrevil/simplified-market-ledger/ledger/transactions"
	"gorm.io/gorm"
)

type PostgresTransactionsRepository struct {
	db *gorm.DB
}

func (p *PostgresTransactionsRepository) GetTransaction(transactionID int32) (transactions.Transaction, error) {
	tr := transactions.Transaction{ID: transactionID}
	err := p.db.First(&tr).Error
	return tr, err
}

func (p *PostgresTransactionsRepository) CreateTransaction(tr transactions.Transaction) (int32, error) {
	result := p.db.Create(&tr)
	return tr.ID, result.Error
}
func (p *PostgresTransactionsRepository) UpdateTransaction(tr transactions.Transaction) error {
	return p.db.Save(&tr).Error
}

func (p *PostgresTransactionsRepository) ListInvestorTransactions(investorID int32) []transactions.Transaction {
	transactionsList := []transactions.Transaction{}
	p.db.Where(transactions.Transaction{
		ID: investorID,
	}).Find(&transactionsList)
	return transactionsList
}
