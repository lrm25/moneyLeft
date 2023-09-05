package models

// BankAccountNoInterest represents an account with no passive increase
type BankAccountNoInterest struct {
	*CashAccount
}

// AccountsNoInterest - simplify for array declaration
type AccountsNoInterest []*BankAccountNoInterest

// NewAccountNoInterest constructor - name, amount in account
func NewAccountNoInterest(name string, amount float64) *BankAccountNoInterest {
	return &BankAccountNoInterest{
		CashAccount: NewCashAccount(name, TypeSavingsNoInterest, amount),
	}
}