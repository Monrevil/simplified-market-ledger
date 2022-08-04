package postgres

import (
	"github.com/Monrevil/simplified-market-ledger/issuers"
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

func (p *PostgresIssuersRepository) ChangeBalance(id int32, amount int) error {
	issuer := &issuers.Issuer{ID: id}
	return p.db.Model(issuer).Update("balance", gorm.Expr("balance + ?", amount)).Error
}
