package models

type PositiveAccount interface {
	Account
	Deduct(amount float64) float64
	String() string
}

type PositiveAccounts []PositiveAccount
