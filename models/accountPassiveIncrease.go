package models

// PassiveIncreaseAccount refers to an account that can passively increase with interest/bond return/etc.
type PassiveIncreaseAccount interface {
	PositiveAccount
	Increase()
}

// PassiveIncreaseAccounts - simplify array declaration
type PassiveIncreaseAccounts []PassiveIncreaseAccount