package investors

type Repository interface {
	GetInvestor(investorID uint) (Investor, error)
	ReserveBalance(investorID uint, amount int) error
	ReleaseBalance(investorID uint, amount int) error
	ChangeReservedBalance(investorID uint, amount int) error
	ListInvestors() []Investor
}

type Investor struct {
	ID              uint
	FirstName       string
	LastName        string
	Balance         int
	ReservedBalance int
}
