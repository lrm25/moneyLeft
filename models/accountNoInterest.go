package models

type AccountNoInterest struct {
	*CashAccount
}

func NewAccountNoInterest(name string, amount float64) *AccountNoInterest {
	return &AccountNoInterest{
		CashAccount: NewCashAccount(name, TypeSavingsNoInterest, amount),
	}
}
