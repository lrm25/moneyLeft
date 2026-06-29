package models

// AccountBond - struct containing bond account info
type AccountBond struct {
	*AccountWithInterest
	person *Person
}

// AccountsBond - simplify array declaration
type AccountsBond []*AccountBond

// NewAccountBond constructor (account name, amount, expected stock interest rate, sale fee, person holding account)
func NewAccountBond(name string, amount, interestRate float64, person *Person) *AccountBond {
	if person == nil {
		panic("Person for stock brokerage account cannot be nil")
	}
	return &AccountBond{
		AccountWithInterest: &AccountWithInterest{
			BankAccount: &BankAccount{
				name:        name,
				amount:      amount,
				accountType: TypeStockBrokerage,
				removable:   true,
			},
			interestRate: interestRate,
		},
		person: person,
	}
}

// Increase the amount for a single month with the yearly interest rate
func (a *AccountBond) Increase() {
	a.amount *= 1 + (a.interestRate / 1200.0)
}

// Person returns the person owning this account
func (a *AccountBond) Person() *Person {
	return a.person
}

// Deduct removes money from the brokerage account when the user retrieves it, taking into account that
// the money is now cap gains taxable and requires a sale fee
func (a *AccountBond) Deduct(amount float64) (float64, float64) {
	a.amount -= amount
	a.person.taxableCapGainsThis += amount
	outstanding := 0.0
	if a.amount <= 0 {
		outstanding = a.amount * -1
		a.closed = true
		a.amount = 0
	}
	return a.amount, outstanding
}
