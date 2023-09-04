package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IRA(t *testing.T) {
	deduct := 1000.00
	person := NewPerson(59, 0, 70, deduct, 0.0)
	iraAccount := NewIRA("test IRA", 100000.00, 50.0, 50.0, 10.0, person)
	stockAmount := 50000.00 * (1 + 50.00/1200.00)
	bondAmount := 50000.00 * (1 + 10.00/1200.00)
	person.WithAccounts(nil, nil, PassiveIncreaseAccounts{iraAccount})
	person.IncreaseAge(2000, 1)
	require.InDelta(t, iraAccount.Amount(), stockAmount+bondAmount-deduct-deduct/10, 1.00)
	require.InDelta(t, person.taxableOtherThis, 1000.00, 1.00)

	deduct = 1000.00
	person = NewPerson(59, 6, 70, deduct, 0.0)
	iraAccount = NewIRA("test IRA", 100000.00, 50.0, 50.0, 10.0, person)
	stockAmount = 50000.00 * (1 + 50.00/1200.00)
	bondAmount = 50000.00 * (1 + 10.00/1200.00)
	person.WithAccounts(nil, nil, PassiveIncreaseAccounts{iraAccount})
	person.IncreaseAge(2000, 1)
	require.InDelta(t, iraAccount.Amount(), stockAmount+bondAmount-deduct, 1.00)
	require.InDelta(t, person.taxableOtherThis, 1000.00, 1.00)
}