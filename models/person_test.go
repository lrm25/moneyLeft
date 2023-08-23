package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PayCreditCardsNothing(t *testing.T) {
	p := NewPerson(50, 0, 80, 1000.00, 0.0)
	ccOne := NewCreditCardAccount("Test cc one", 1000.00)
	p.WithAccounts([]*CreditCardAccountImpl{ccOne}, nil, nil)
	p.PayCreditCards()
	require.True(t, p.broke, "Person should be broke with credit card but no accounts")

	p = NewPerson(50, 0, 80, 1000.00, 0.0)
	baOne := NewAccountNoInterest("test account one", 500.00)
	p.WithAccounts([]*CreditCardAccountImpl{ccOne}, []PositiveAccount{baOne}, nil)
	p.PayCreditCards()
	require.True(t, p.broke, "Person should be broke")
	require.InDelta(t, ccOne.Amount(), 500.00, 0.01, "500.00 should be paid")

	p = NewPerson(50, 0, 80, 1000.00, 0.0)
	ccOne = NewCreditCardAccount("Test cc one", 1000.00)
	baOne = NewAccountNoInterest("ba one", 600.00)
	baTwo := NewAccountWithInterest("ba two", 300.00, 0.10)
	p.WithAccounts([]*CreditCardAccountImpl{ccOne}, PositiveAccounts{baOne}, PassiveIncreaseAccounts{baTwo})
	p.PayCreditCards()
	require.True(t, p.broke, "Person should be broke")
	require.InDelta(t, ccOne.Amount(), 100.00, 0.01, "100.00 should be paid")

	p = NewPerson(50, 0, 80, 1000.00, 0.0)
	ccOne = NewCreditCardAccount("Test cc one", 750.00)
	baOne = NewAccountNoInterest("ba one", 600.00)
	baTwo = NewAccountWithInterest("ba two", 300.00, 0.10)
	p.WithAccounts([]*CreditCardAccountImpl{ccOne}, PositiveAccounts{baOne}, PassiveIncreaseAccounts{baTwo})
	p.PayCreditCards()
	require.True(t, ccOne.Closed(), "Credit card account should be closed")
	require.False(t, p.broke, "Person should not be broke")
	require.InDelta(t, baTwo.Amount(), 150.00, 0.01, "150.00 should be remaining in account")

	p = NewPerson(50, 0, 80, 1000.00, 0.0)
	ccOne = NewCreditCardAccount("Test cc one", 750.00)
	ccTwo := NewCreditCardAccount("Test cc two", 500.00)
	baOne = NewAccountNoInterest("ba one", 600.00)
	baTwo = NewAccountWithInterest("ba two", 300.00, 0.10)
	p.WithAccounts([]*CreditCardAccountImpl{ccOne, ccTwo}, PositiveAccounts{baOne}, PassiveIncreaseAccounts{baTwo})
	p.PayCreditCards()
	require.True(t, ccOne.Closed(), "Credit card account one should be closed")
	require.False(t, ccTwo.Closed(), "Credit card account two should not be closed")
	require.True(t, p.broke, "Person should be broke")
	require.InDelta(t, ccTwo.Amount(), 350.00, 0.01, "350.00 should be remaining on credit card two")
}

func Test_SocialSecurity(t *testing.T) {
	p := NewPerson(61, 0, 70, 0.0, 0.0)
	payouts := NewPayoutsSocialSecurity(1000.00, 2000.00, 3000.00)
	ssAccount := NewAccountSocialSecurity(Early, payouts, 10.0, p)
	accounts := PassiveIncreaseAccounts{}
	accounts = append(accounts, ssAccount)
	p.WithAccounts(nil, nil, accounts)
	for idx := 0; idx < 12; idx++ {
		p.IncreaseAge(2000, 1)
	}
	assert.InDelta(t, ssAccount.Amount(), 1100.00, 10.00)
	p.IncreaseAge(2001, 1)
	assert.InDelta(t, ssAccount.Amount(), 2200.00, 20.00)
	p.IncreaseAge(2001, 1)
	assert.InDelta(t, ssAccount.Amount(), 3200.00, 20.00)

	p = NewPerson(65, 0, 70, 0.0, 0.0)
	ssAccount = NewAccountSocialSecurity(Normal, payouts, 10.0, p)
	accounts = PassiveIncreaseAccounts{}
	accounts = append(accounts, ssAccount)
	p.WithAccounts(nil, nil, accounts)
	for idx := 0; idx < 24; idx++ {
		p.IncreaseAge(2000, 1)
	}
	assert.InDelta(t, ssAccount.Amount(), 2420.00, 100.00)
	p.IncreaseAge(2002, 1)
	assert.InDelta(t, ssAccount.Amount(), 4900.00, 20.00)
}

func Test_IncreaseAge(t *testing.T) {
	person := NewPerson(40, 0, 50, 1000.00, 1.0)
	bankAccount := NewAccountNoInterest("test no interest", 1000.00)
	investmentAccount := NewAccountWithInterest("test with interest", 1000.00, 10)
	println("amount", investmentAccount.Amount())
	person.WithAccounts(nil, PositiveAccounts{bankAccount}, PassiveIncreaseAccounts{investmentAccount})
	person.IncreaseAge(2000, 1)
	assert.False(t, person.broke)
	person.IncreaseAge(2000, 2)
	assert.False(t, person.broke)
	person.IncreaseAge(2000, 3)
	assert.True(t, person.broke)
}

func Test_PayTaxes(t *testing.T) {
	person := NewPerson(40, 0, 50, 1000.00, 1.0)
	bankAccount := NewAccountNoInterest("test no interest", 1000.00)
	investmentAccount := NewAccountWithInterest("test interest", 10000.00, 10)
	person.taxableOtherThis = 5000
	person.taxableCapGainsThis = 5000
	person.nonCapBrackets = NewTaxBrackets(4000.00, []*TaxBracket{NewTaxBracket(0.0, 4000.00, 10.0)})
	person.capBrackets = NewCapTaxBrackets([]*TaxBracket{NewTaxBracket(0.0, 20000.0, 20.0)})
	person.WithAccounts(nil, PositiveAccounts{bankAccount}, PassiveIncreaseAccounts{investmentAccount})
	person.PayTaxes()
	assert.InDelta(t, 9900.00, investmentAccount.Amount(), 1.00)
}
