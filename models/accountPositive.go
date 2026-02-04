package models

// PositiveAccount refers to an account with a positive amount of money, or non-credit card account
type PositiveAccount interface {
	Account
	Deduct(amount float64) (float64, float64)
	String() string
}

// PositiveAccounts - simplify array declaration
type PositiveAccounts []PositiveAccount
