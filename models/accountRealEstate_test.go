package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRealEstate_Deduct(t *testing.T) {
	person := NewPerson(30, 0, 70, 1000.00, 0.00)
	testRealEstate := NewRealEstate(
		&RealEstateConfig{
			person,
			"realEstate",
			100000,
			0.00,
			5,
			5,
			40,
			50,
			12,
			10,
			6,
		})
	remainingInAccount, remainingToPay := testRealEstate.Deduct(1000.00)
	require.InDelta(t, remainingInAccount, 0.00, 0.01)
	require.InDelta(t, remainingToPay, 1000.00, 0.01)
	require.Equal(t, len(person.accounts), 1)
	require.Equal(t, len(person.interestAccounts), 2)
	cashAccount := person.accounts[0]
	require.Equal(t, "realEstate sale cash", cashAccount.Name())
	require.InDelta(t, 5000.00, cashAccount.Amount(), 0.01)

	stockAccountChecked := false
	bondAccountChecked := false
	for _, account := range person.interestAccounts {
		if account.Name() == "realEstate sale stock" {
			require.Equal(t, 50000.00, account.Amount(), 0.01)
			require.Equal(t, 12.00, account.(*AccountStockBrokerage).interestRate)
			require.Equal(t, 10.00, account.(*AccountStockBrokerage).MonthlySaleFee())
			stockAccountChecked = true
		} else if account.Name() == "realEstate sale bonds" {
			require.Equal(t, 40000.00, account.Amount(), 0.01)
			require.Equal(t, 6.00, account.(*AccountStockBrokerage).interestRate)
			bondAccountChecked = true
		}
	}
	require.True(t, stockAccountChecked)
	require.True(t, bondAccountChecked)
}
