package models

import (
	"fmt"
	"github.com/lrm25/moneyLeft/logger"
)

// IRA - structure representing IRA
type IRA struct {
	*BankAccount
	percentStock      float64
	stockInterestRate float64
	bondInterestRate  float64
	person            *Person
}

// IRAs - simplify array declaration
type IRAs []*IRA

// NewIRA constructor (account name, amount, percent that is stock as opposed to bond, stock and bond expected interest rates,
// person holding the account)
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

// PercentStock returns the percentage of the IRA that is stock.  As of now, it is assumed that the rest is in bonds.
func (i *IRA) PercentStock() float64 {
	return i.percentStock
}

// StockInterestRate returns the expected stock interest rates
func (i *IRA) StockInterestRate() float64 {
	return i.stockInterestRate
}

// BondInterestRate returns the expected bond return rates for the bonds in the IRA
func (i *IRA) BondInterestRate() float64 {
	return i.bondInterestRate
}

// Person returns who holds this account
func (i *IRA) Person() *Person {
	return i.person
}

// Increase (passively) the amount of money over a month using the expected yearly and monthly returns
func (i *IRA) Increase() {
	totalStock := i.amount * (i.percentStock / 100.0)
	totalBond := i.amount - totalStock
	totalStock *= 1 + (i.stockInterestRate / 1200.0)
	totalBond *= 1 + (i.bondInterestRate / 1200.0)
	i.amount = totalStock + totalBond
	logger.Get().Debug(fmt.Sprintf("Increasing %s amount to %.2f", i.name, i.amount))
}

// Deduct from the account.  Take into account taxable income, and if the holder is less than 59.5 years old,
// deduct a 10% penalty as well
func (i *IRA) Deduct(amount float64) (float64, float64) {
	i.amount -= amount
	logger.Get().Debug(fmt.Sprintf("Deducting %.2f from %s", amount, i.name))
	i.person.taxableOtherThis += amount
	if i.person.years < 59 || (i.person.years == 59 && i.person.months < 6) {
		penalty := amount * 0.1
		i.amount -= penalty
		logger.Get().Debug(fmt.Sprintf("Deducting %.2f penalty from %s", penalty, i.name))
	}
	outstanding := 0.0
	if i.amount <= 0 {
		i.closed = true
		outstanding = -1 * i.amount
		i.amount = 0
	}
	return i.amount, outstanding
}
