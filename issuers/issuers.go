package issuers

type Repository interface {
	ChangeBalance(id uint, amount int)
}

type Issuer struct {
	ID        uint
	FirstName string
	LastName  string
	Balance   int
}
