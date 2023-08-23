package models

type AccountStockBrokerage struct {
	*AccountWithInterest
	BrokerageFeePerMonth float64
	Person               *Person
}

func NewAccountStockBrokerage(name string, amount, interestRate, brokerageFee float64, person *Person) *AccountStockBrokerage {
	return &AccountStockBrokerage{
		AccountWithInterest: &AccountWithInterest{
			BankAccount: &BankAccount{
				name:        name,
				amount:      amount,
				accountType: TypeStockBrokerage,
				removable:   true,
			},
			InterestRate: interestRate,
		},
		BrokerageFeePerMonth: brokerageFee,
		Person:               person,
	}
}

func (a *AccountStockBrokerage) Increase() {
	a.amount *= (1 + (a.InterestRate / 1200.0))
}

func (a *AccountStockBrokerage) Closed() bool {
	return a.closed
}

func (a *AccountStockBrokerage) Deduct(amount float64) float64 {
	a.amount -= amount
	a.amount -= a.BrokerageFeePerMonth
	a.Person.taxableCapGainsThis += amount
	if a.amount <= 0 {
		a.closed = true
	}
	return a.amount
}
