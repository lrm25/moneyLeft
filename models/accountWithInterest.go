package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

type AccountWithInterest struct {
	*BankAccount
	interestRate float64
}

type AccountsWithInterest []*AccountWithInterest

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

func (a *AccountWithInterest) Rate() float64 {
	return a.interestRate
}

func (a *AccountWithInterest) Increase() {
	logger.Get().Debug(fmt.Sprintf("increasing %s: %.2f\n", a.name, a.amount))
	a.amount *= (1 + (a.interestRate / 1200.0))
}