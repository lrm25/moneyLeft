package models

import "fmt"

// NOTE:  I pay my cards off every month, if not this would have to change
type CreditCardAccountImpl BankAccount

func NewCreditCardAccount(name string, debt float64) *CreditCardAccountImpl {
	return &CreditCardAccountImpl{
		name:        name,
		amount:      debt,
		accountType: TypeCreditCard,
		removable:   true,
	}
}

func (c *CreditCardAccountImpl) Type() int {
	return TypeCreditCard
}

func (c *CreditCardAccountImpl) Amount() float64 {
	return c.amount
}

func (c *CreditCardAccountImpl) Closed() bool {
	return c.closed
}

func (c *CreditCardAccountImpl) Close() {
	c.amount = 0
	c.closed = true
}

func (c *CreditCardAccountImpl) Pay(account PositiveAccount) bool {
	accountAmount := account.Amount()
	if c.amount <= accountAmount {
		account.Deduct(c.amount)
		c.Close()
	} else {
		account.Close()
		c.amount -= accountAmount
	}
	fmt.Printf("Credit card amount for %s: %.2f", c.name, c.amount)
	return c.closed
}
