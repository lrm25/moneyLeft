package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

// CreditCards - simplify array declaration
type CreditCards []*CreditCardAccountImpl

// CreditCardAccountImpl represents a credit card account
type CreditCardAccountImpl BankAccount

// NewCreditCardAccount constructor (name, remaining to pay off)
func NewCreditCardAccount(name string, debt float64) *CreditCardAccountImpl {
	return &CreditCardAccountImpl{
		name:        name,
		amount:      debt,
		accountType: TypeCreditCard,
		removable:   true,
	}
}

// Name returns the credit card name
func (c *CreditCardAccountImpl) Name() string {
	return c.name
}

// Type returns the account type, in this case, credit card
func (c *CreditCardAccountImpl) Type() int {
	return TypeCreditCard
}

// Amount returns the amount remaining to pay off
func (c *CreditCardAccountImpl) Amount() float64 {
	return c.amount
}

// Closed returns whether or not the credit card has been paid off
func (c *CreditCardAccountImpl) Closed() bool {
	return c.closed
}

// Close the credit card account, if paid off
func (c *CreditCardAccountImpl) Close() {
	c.amount = 0
	c.closed = true
}

// Pay the credit card.  This program assumes the user pays the credit cards off each month and does not take credit card interest
// into account.
func (c *CreditCardAccountImpl) Pay(account PositiveAccount) bool {
	logger.Get().Debug(fmt.Sprintf("Credit card amount for %s before payment: %.2f", c.name, c.amount))
	logger.Get().Debug(fmt.Sprintf("Paying from %s with total amount: %.2f", account.Name(), account.Amount()))
	//accountAmount := account.Amount()
	_, c.amount = account.Deduct(c.amount)
	if c.amount < 0.001 {
		c.Close()
	}
	logger.Get().Debug(fmt.Sprintf("Credit card amount for %s after payment: %.2f", c.name, c.amount))
	return c.closed
}
