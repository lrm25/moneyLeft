package models

type AccountStockBrokerage struct {
	*AccountWithInterest
	brokerageFeePerMonth float64
	person               *Person
}

type AccountsStockBrokerage []*AccountStockBrokerage

func NewAccountStockBrokerage(name string, amount, interestRate, brokerageFee float64, person *Person) *AccountStockBrokerage {
	return &AccountStockBrokerage{
		AccountWithInterest: &AccountWithInterest{
			BankAccount: &BankAccount{
				name:        name,
				amount:      amount,
				accountType: TypeStockBrokerage,
				removable:   true,
			},
			interestRate: interestRate,
		},
		brokerageFeePerMonth: brokerageFee,
		person:               person,
	}
}

func (a *AccountStockBrokerage) Increase() {
	a.amount *= (1 + (a.interestRate / 1200.0))
}

func (a *AccountStockBrokerage) Person() *Person {
	return a.person
}

func (a *AccountStockBrokerage) MonthlySaleFee() float64 {
	return a.brokerageFeePerMonth
}

func (a *AccountStockBrokerage) Deduct(amount float64) float64 {
	a.amount -= amount
	a.amount -= a.brokerageFeePerMonth
	a.person.taxableCapGainsThis += amount
	if a.amount <= 0 {
		a.closed = true
	}
	return a.amount
}