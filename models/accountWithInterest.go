package models

import "fmt"

type AccountWithInterest struct {
	*BankAccount
	InterestRate float64
}

func NewAccountWithInterest(name string, amount, interestRate float64) *AccountWithInterest {
	return &AccountWithInterest{
		BankAccount: &BankAccount{
			name:        name,
			amount:      amount,
			accountType: TypeSavingsWithInterest,
			removable:   true,
		},
		InterestRate: interestRate,
	}
}

func (a *AccountWithInterest) Increase() {
	fmt.Printf("increasing %s: %.2f\n", a.name, a.amount)
	a.amount *= (1 + (a.InterestRate / 1200.0))
}
