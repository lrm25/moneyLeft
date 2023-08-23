package models

type IRA struct {
	*BankAccount
	PercentStock      float64
	StockInterestRate float64
	BondInterestRate  float64
	person            *Person
}

func NewIRA(name string, amount, percentStock, stockInterestRate, bondInterestRate float64, person *Person) *IRA {
	return &IRA{
		BankAccount: &BankAccount{
			name:        name,
			amount:      amount,
			accountType: TypeIRA,
			removable:   true,
		},
		PercentStock:      percentStock,
		StockInterestRate: stockInterestRate,
		BondInterestRate:  bondInterestRate,
		person:            person,
	}
}

func (i *IRA) Increase() {
	totalStock := i.amount * (i.PercentStock / 100.0)
	totalBond := i.amount - totalStock
	totalStock *= (1 + (i.StockInterestRate / 1200.0))
	totalBond *= (1 + (i.BondInterestRate / 1200.0))
	i.amount = totalStock + totalBond
}

func (i *IRA) Deduct(amount float64) float64 {
	i.amount -= amount
	i.person.taxableOtherThis += amount
	if i.person.years < 60 && i.person.months < 6 {
		i.amount -= amount * 0.1
	}
	if i.amount < 0 {
		i.closed = true
	}
	return i.amount
}
