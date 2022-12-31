package account

type AccountType int

const (
	SavingsNoInterest AccountType = iota + 1
	SavingsWithInterest
	StockBrokerage
	FourOhOneK
	SocialSecurity
)

type Account struct {
	name        string
	accountType AccountType
	amount      float64
}
