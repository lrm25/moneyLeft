package models

import "fmt"

type SocialSecuritySelection int

const (
	Early  SocialSecuritySelection = 62
	Normal                         = 67
	Late                           = 70
)

type PayoutsSocialSecurity struct {
	early  float64
	normal float64
	late   float64
}

type AccountSocialSecurity struct {
	*AccountWithInterest
	selection SocialSecuritySelection
	perMonth  float64
	Person    *Person
}

func NewPayoutsSocialSecurity(early, normal, late float64) *PayoutsSocialSecurity {
	return &PayoutsSocialSecurity{
		early:  early,
		normal: normal,
		late:   late,
	}
}

func NewAccountSocialSecurity(selection SocialSecuritySelection, payouts *PayoutsSocialSecurity, interestRate float64, person *Person) *AccountSocialSecurity {
	var payout float64
	if selection == Early {
		payout = payouts.early
	} else if selection == Normal {
		payout = payouts.normal
	} else if selection == Late {
		payout = payouts.late
	} else {
		panic(fmt.Sprintf("Invalid payout code: %v", payouts))
	}
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
			InterestRate: interestRate,
		},
		Person: person,
	}
}

func (s *AccountSocialSecurity) WithPerson(person *Person) *AccountSocialSecurity {
	s.Person = person
	return s
}

func (s *AccountSocialSecurity) Closed() bool {
	println("SS", s.Person.years, s.selection)
	if s.Person.years < int(s.selection) {
		s.closed = true
		return true
	}
	s.closed = false
	return false
}

func (a *AccountSocialSecurity) Increase() {
	a.perMonth *= (1 + (a.InterestRate / 1200.00))
	if !a.Closed() {
		a.amount += a.perMonth
		fmt.Printf("increasing social security to %.2f\n", a.amount)
	}
}

func (a *AccountSocialSecurity) Amount() float64 {
	if a.Person.years < int(a.selection) {
		return 0
	}
	return a.amount
}

func (a *AccountSocialSecurity) Deduct(amount float64) float64 {
	fmt.Printf("deducting: %.2f\n", amount)
	a.amount -= amount
	if a.amount <= 0 {
		a.closed = true
	}
	return a.amount
}
