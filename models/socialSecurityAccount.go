package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

type SocialSecuritySelection int

const (
	EARLY SocialSecuritySelection  = 62
	NORMAL SocialSecuritySelection = 67
	LATE SocialSecuritySelection   = 70
)

type AccountSocialSecurity struct {
	*AccountWithInterest
	selection SocialSecuritySelection
	perMonth  float64
	person    *Person
}

func NewAccountSocialSecurity(selection SocialSecuritySelection, payout float64, interestRate float64, person *Person) *AccountSocialSecurity {
	return &AccountSocialSecurity{
		selection: selection,
		perMonth:  payout,
		AccountWithInterest: &AccountWithInterest{
			BankAccount: &BankAccount{
				name:        "Social Security",
				amount:      payout,
				accountType: TypeSocialSecurity,
				removable:   false,
			},
			interestRate: interestRate,
		},
		person: person,
	}
}

func (s *AccountSocialSecurity) Person() *Person {
	return s.person
}

func (s *AccountSocialSecurity) Selection() int {
	return int(s.selection)
}

func (s *AccountSocialSecurity) Closed() bool {
	logger.Get().Debug(fmt.Sprintf("SS:  years %d, selection %d", s.person.years, s.selection))
	if s.person.years < int(s.selection) {
		s.closed = true
		return true
	}
	s.closed = false
	return false
}

func (a *AccountSocialSecurity) Increase() {
	a.perMonth *= (1 + (a.interestRate / 1200.00))
	// don't increase on first month available
	if !a.Closed() && (int(a.selection) < a.person.years || 0 < a.person.months) {
		a.amount += a.perMonth
		logger.Get().Debug(fmt.Sprintf("increasing social security to %.2f\n", a.amount))
		return
	}
	a.amount = a.perMonth
}

func (a *AccountSocialSecurity) Amount() float64 {
	if a.person.years < int(a.selection) {
		return 0
	}
	return a.amount
}