package postgres

import (
	"github.com/Monrevil/simplified-market-ledger/ledger/issuers"
	"gorm.io/gorm"
)

// Issuers:

type PostgresIssuersRepository struct {
	db *gorm.DB
}

func (p *PostgresIssuersRepository) NewIssuer(iss *issuers.Issuer) error {
	return p.db.Create(iss).Error
}

func (p *PostgresIssuersRepository) GetIssuer(iss *issuers.Issuer) error {
	return p.db.First(iss).Error
}

func (p *PostgresIssuersRepository) IncreaseBalance(id int32, amount int32) error {
	issuer := &issuers.Issuer{ID: id}
	return p.db.Model(issuer).Update("balance", gorm.Expr("balance + ?", amount)).Error
}
