package models

// AccountStockBrokerage - struct containing stock brokerage account info
type AccountStockBrokerage struct {
	*AccountWithInterest
	brokerageFeePerMonth float64
	person               *Person
}

// AccountsStockBrokerage - simplify array declaration
type AccountsStockBrokerage []*AccountStockBrokerage

// NewAccountStockBrokerage constructor (accout name, amount, expected stock interest rate, sale fee, person holding account)
func NewAccountStockBrokerage(name string, amount, interestRate, brokerageFee float64, person *Person) *AccountStockBrokerage {
	if person == nil {
		panic("Person for stock brokerage account cannot be nil")
	}
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

// Increase the amount for a single month with the yearly interest rate
func (a *AccountStockBrokerage) Increase() {
	a.amount *= 1 + (a.interestRate / 1200.0)
}

// Person returns the person owning this account
func (a *AccountStockBrokerage) Person() *Person {
	return a.person
}

// MonthlySaleFee returns the estimated fee for selling stocks over a month period
func (a *AccountStockBrokerage) MonthlySaleFee() float64 {
	return a.brokerageFeePerMonth
}

// Deduct removes money from the brokerage account when the user retrieves it, taking into account that
// the money is now cap gains taxable and requires a sale fee
func (a *AccountStockBrokerage) Deduct(amount float64) (float64, float64) {
	a.amount -= amount
	a.amount -= a.brokerageFeePerMonth
	a.person.taxableCapGainsThis += amount
	outstanding := 0.0
	if a.amount <= 0 {
		outstanding = a.amount * -1
		a.closed = true
		a.amount = 0
	}
	return a.amount, outstanding
}
