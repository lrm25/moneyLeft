package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

// AccountWithInterest represents an account with a single expected increasing interest rate (bank accounts w/interest, stock brokerage accounts, etc.)
type AccountWithInterest struct {
	*BankAccount
	interestRate float64
}

// AccountsWithInterest - simplify array declaration
type AccountsWithInterest []*AccountWithInterest

// NewAccountWithInterest constructor (name, amount, interest rate)
func NewAccountWithInterest(name string, amount, interestRate float64) *AccountWithInterest {
	return &AccountWithInterest{
		BankAccount: &BankAccount{
			name:        name,
			amount:      amount,
			accountType: TypeSavingsWithInterest,
			removable:   true,
		},
		interestRate: interestRate,
	}
}

// Rate returns the interest rate for this account
func (a *AccountWithInterest) Rate() float64 {
	return a.interestRate
}

// Increase adds the passive increase for this account
func (a *AccountWithInterest) Increase() {
	logger.Get().Debug(fmt.Sprintf("increasing %s: %.2f\n", a.name, a.amount))
	a.amount *= (1 + (a.interestRate / 1200.0))
}