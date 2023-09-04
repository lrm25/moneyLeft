package models

type BankAccountNoInterest struct {
	*CashAccount
}

type AccountsNoInterest []*BankAccountNoInterest

func NewAccountNoInterest(name string, amount float64) *BankAccountNoInterest {
	return &BankAccountNoInterest{
		CashAccount: NewCashAccount(name, TypeSavingsNoInterest, amount),
	}
}