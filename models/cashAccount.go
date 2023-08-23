package models

type CashAccount struct {
	*BankAccount
}

func NewCashAccount(name string, typ int, amount float64) *CashAccount {
	return &CashAccount{
		BankAccount: NewBankAccount(name, typ, amount),
	}
}

func (a *CashAccount) Deduct(amount float64) float64 {
	result := a.amount - amount
	if result <= 0 {
		a.amount = 0.0
		if a.accountType != TypeSocialSecurity {
			a.closed = true
		}
	} else {
		a.amount = result
	}
	return result
}
