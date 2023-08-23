package models

import "fmt"

const (
	TypeCreditCard = iota + 1
	TypeSavingsNoInterest
	TypeSavingsWithInterest
	TypeStockBrokerage
	TypeIRA
	TypeSocialSecurity
)

type Account interface {
	Amount() float64
	Name() string
	Type() int
	Closed() bool
	Removable() bool
	Close()
}

type BankAccount struct {
	name        string
	accountType int
	amount      float64
	closed      bool
	removable   bool
}

func (b *BankAccount) String() string {
	return fmt.Sprintf("Name: %s, accountType: %d, amount: %.2f, closed: %t, removable: %t", b.name, b.accountType, b.amount, b.closed, b.removable)
}

func NewBankAccount(name string, typ int, amount float64) *BankAccount {
	return &BankAccount{
		name:        name,
		accountType: typ,
		amount:      amount,
		removable:   true,
	}
}

func (a *BankAccount) Amount() float64 {
	return a.amount
}

func (a *BankAccount) Name() string {
	return a.name
}

func (a *BankAccount) Type() int {
	return a.accountType
}

func (a *BankAccount) Closed() bool {
	return a.closed
}

func (a *BankAccount) Close() {
	a.amount = 0
	a.closed = true
}

func (a *BankAccount) Deduct(amount float64) float64 {
	fmt.Printf("deducting: %.2f\n", amount)
	a.amount -= amount
	if a.amount <= 0 {
		a.closed = true
	}
	return a.amount
}

func (a *BankAccount) Removable() bool {
	return a.removable
}

func (a Accounts) Len() int {
	return len(a)
}

func (a Accounts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Accounts) Less(i, j int) bool {
	if a[i].Type() < a[j].Type() {
		return true
	}
	return a[i].Amount() < a[j].Amount()
}

type Accounts []Account
