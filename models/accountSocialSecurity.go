package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

// SocialSecuritySelection represents the age that the user selects to receive social security
type SocialSecuritySelection int

// Selection ages
const (
	Early  SocialSecuritySelection = 62
	Normal SocialSecuritySelection = 67
	Late   SocialSecuritySelection = 70
)

// AccountSocialSecurity - struct
type AccountSocialSecurity struct {
	*AccountWithInterest
	selection         SocialSecuritySelection
	perMonth          float64
	person            *Person
	yearlyDropPercent int
}

// NewAccountSocialSecurity constructor (selection - early, normal, or late, monthly payout for selection as of now, expected inflation rate, person with account)
func NewAccountSocialSecurity(selection SocialSecuritySelection, payout float64, interestRate float64, person *Person, yearlyDropPercent int) *AccountSocialSecurity {
	return &AccountSocialSecurity{
		selection:         selection,
		perMonth:          payout,
		yearlyDropPercent: yearlyDropPercent,
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

// Person returns person with social security account
func (s *AccountSocialSecurity) Person() *Person {
	return s.person
}

// Selection returns user's social security selection (early - 62 years, normal - 67 years, late - 70 years)
func (s *AccountSocialSecurity) Selection() int {
	return int(s.selection)
}

// Closed is closed if the user hasn't reached age yet.
func (s *AccountSocialSecurity) Closed() bool {
	logger.Get().Debug(fmt.Sprintf("SS:  years %d, selection %d", s.person.years, s.selection))
	if s.person.years < int(s.selection) {
		s.closed = true
		return true
	}
	s.closed = false
	return false
}

// Increase the amount in the account.  Before the user is old enough to receive, just adjust future per-month payment based on inflation.
// If the user is old enough, add amount to account each month.
func (s *AccountSocialSecurity) Increase() {
	s.perMonth *= 1 + (s.interestRate / 1200.00)
	logger.Get().Debug(fmt.Sprintf("increasing social security monthly amount to %.2f\n", s.perMonth))
	// don't increase on first month available
	if !s.Closed() && (int(s.selection) < s.person.years || 0 < s.person.months) {
		s.amount += s.perMonth
		logger.Get().Debug(fmt.Sprintf("increasing social security to %.2f\n", s.amount))
		return
	}
	// reduce (assumption is that you're not getting paid)
	if s.Closed() {
		s.perMonth *= 1 - (float64(s.yearlyDropPercent) / 1200.00)
		logger.Get().Debug(fmt.Sprintf("decreasing future social security to %.2f\n", s.perMonth))
	}
	s.amount = s.perMonth
}

// Amount deductible from account.  Should be zero unless the user's age is high enough to receive social security.
func (s *AccountSocialSecurity) Amount() float64 {
	if s.person.years < int(s.selection) {
		return 0
	}
	return s.amount
}
