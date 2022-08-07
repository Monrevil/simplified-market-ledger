package postgres

import (
	"fmt"

	"github.com/Monrevil/simplified-market-ledger/ledger/investors"
	"gorm.io/gorm"
)

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
func (p *PostgresInvestorsRepository) ReserveBalance(investor *investors.Investor, amount int32) error {
	if investor.ID == 0 {
		return fmt.Errorf("InvestorID can not be 0")
	}
	if err := p.db.First(investor).Error; err != nil {
		return err
	}

	if investor.Balance < amount {
		return fmt.Errorf("not enough funds. Tried to reserve %v funds. With %v active balance", amount, investor.Balance)
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
func (p *PostgresInvestorsRepository) ReleaseBalance(investor *investors.Investor, amount int32) error {
	if err := p.db.First(investor).Error; err != nil {
		return err
	}
	if investor.ReservedBalance < amount {
		return fmt.Errorf("not enough funds. Tried to release %v funds. With %v reserved balance", amount, investor.ReservedBalance)
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

func (p *PostgresInvestorsRepository) ReduceReservedBalance(investor *investors.Investor, amount int32) error {
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
