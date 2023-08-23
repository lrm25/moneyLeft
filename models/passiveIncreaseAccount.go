package models

type PassiveIncreaseAccount interface {
	PositiveAccount
	Increase()
}

type PassiveIncreaseAccounts []PassiveIncreaseAccount
