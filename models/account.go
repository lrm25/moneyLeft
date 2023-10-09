package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

// Account types, for sorting
const (
	TypeCreditCard = iota + 1
	TypeSavingsNoInterest
	TypeSavingsWithInterest
	TypeStockBrokerage
	TypeIRA
	TypeSocialSecurity
)

// Account represents any bank account with positive or negative amount
type Account interface {
	Amount() float64
	Name() string
	Type() int
	Closed() bool
	Removable() bool
	Close()
}

// Accounts - simplify array declaration
type Accounts []Account

// BankAccount represents any account with a positive or negative amount of money
type BankAccount struct {
	name        string
	accountType int
	amount      float64
	closed      bool
	removable   bool
}

// String represents a bank account string
func (a *BankAccount) String() string {
	return fmt.Sprintf("Name: %s, accountType: %d, amount: %.2f, closed: %t, removable: %t", a.name, a.accountType, a.amount, a.closed, a.removable)
}

// NewBankAccount constructor - account name, type (see the constants above), and initial amount
func NewBankAccount(name string, typ int, amount float64) *BankAccount {
	return &BankAccount{
		name:        name,
		accountType: typ,
		amount:      amount,
		removable:   true,
	}
}

// Amount in account
func (a *BankAccount) Amount() float64 {
	return a.amount
}

// Name - account name
func (a *BankAccount) Name() string {
	return a.name
}

// Type - credit card, non-interest, IRA, etc.
func (a *BankAccount) Type() int {
	return a.accountType
}

// Closed returns whether or not this account is closed.  All accounts, with the exception of social security, should be set to this
// if the amount reaches zero.
func (a *BankAccount) Closed() bool {
	return a.closed
}

// Close the account
func (a *BankAccount) Close() {
	a.amount = 0
	a.closed = true
}

// Deduct from the account.  If amount is below zero, close.  Return remaining amount in account.
func (a *BankAccount) Deduct(amount float64) float64 {
	logger.Get().Debug(fmt.Sprintf("deducting: %.2f\n", amount))
	a.amount -= amount
	if a.amount <= 0 {
		a.closed = true
	}
	return a.amount
}

// Removable specifies whether or not an account can be removed from the user's account list if the amount reaches zero.
// All accounts except social security should be this.
func (a *BankAccount) Removable() bool {
	return a.removable
}

// Len - for account sorting
func (a Accounts) Len() int {
	return len(a)
}

// Swap - for account sorting
func (a Accounts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less - for account sorting
func (a Accounts) Less(i, j int) bool {
	if a[i].Type() < a[j].Type() {
		return true
	}
	return a[i].Amount() < a[j].Amount()
}