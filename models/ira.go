package models

type IRA struct {
	*BankAccount
	percentStock      float64
	stockInterestRate float64
	bondInterestRate  float64
	person            *Person
}

type IRAs []*IRA

func NewIRA(name string, amount, percentStock, stockInterestRate, bondInterestRate float64, person *Person) *IRA {
	return &IRA{
		BankAccount: &BankAccount{
			name:        name,
			amount:      amount,
			accountType: TypeIRA,
			removable:   true,
		},
		percentStock:      percentStock,
		stockInterestRate: stockInterestRate,
		bondInterestRate:  bondInterestRate,
		person:            person,
	}
}

func (i *IRA) PercentStock() float64 {
	return i.percentStock
}

func (i *IRA) StockInterestRate() float64 {
	return i.stockInterestRate
}

func (i *IRA) BondInterestRate() float64 {
	return i.bondInterestRate
}

func (i *IRA) Person() *Person {
	return i.person
}

func (i *IRA) Increase() {
	totalStock := i.amount * (i.percentStock / 100.0)
	totalBond := i.amount - totalStock
	totalStock *= (1 + (i.stockInterestRate / 1200.0))
	totalBond *= (1 + (i.bondInterestRate / 1200.0))
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