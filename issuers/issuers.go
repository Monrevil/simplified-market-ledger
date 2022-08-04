package issuers

type Repository interface {
	ChangeBalance(id uint, amount int) error
}

type Issuer struct {
	ID        int32
	FirstName string
	LastName  string
	Balance   int
}
