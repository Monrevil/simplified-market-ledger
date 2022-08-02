package investors

type Repository interface {
	GetInvestor(uint) (Investor, error)
	// (InvestorID, Amount)
	ReserveBalance(uint, int) error
	ReleaseBalance(uint, int) error
	ChangeReservedBalance(uint, int) error
	ListInvestors() []Investor
}

type Investor struct {
	ID              uint
	FirstName       string
	LastName        string
	Balance         int
	ReservedBalance int
}
