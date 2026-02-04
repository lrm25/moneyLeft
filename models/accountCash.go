package models

// CashAccount refers to a bank account which contains cash, i.e. not stocks/bonds/etc.
type CashAccount struct {
	*BankAccount
}

// NewCashAccount constructor:  name, type (interest or not (or very tiny and I've disregarded it)), starting amount
func NewCashAccount(name string, typ int, amount float64) *CashAccount {
	return &CashAccount{
		BankAccount: NewBankAccount(name, typ, amount),
	}
}

// Deduct removes a certain amount of cash from the account, and closes it if the amount is below 0.
// The remainder is returned, which is 0 or negative if the account is closed.
func (a *CashAccount) Deduct(amount float64) (float64, float64) {
	result := a.amount - amount
	if result <= 0 {
		a.amount = 0.0
		result *= -1
		if a.accountType != TypeSocialSecurity {
			a.closed = true
		}
	} else {
		a.amount = result
		result = 0
	}
	return a.amount, result
}
