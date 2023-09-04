package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RemovableAccounts(t *testing.T) {

	var accounts Accounts

	nonSSAccount := NewBankAccount("test non ss account", TypeSavingsNoInterest, 0.00)
	ssAccount := NewAccountSocialSecurity(EARLY, 0.00, 0.00, NewPerson(50, 0, 80, 1000.00, 0.0))
	accounts = append(accounts, nonSSAccount, ssAccount)
	for idx, account := range accounts {
		if account.Removable() {
			accounts = append(accounts[:idx], accounts[idx+1:]...)
		}
	}
	require.Equal(t, 1, len(accounts), "One account should be left")
	require.Equal(t, TypeSocialSecurity, accounts[0].Type(), "Social security account should be remaining")
}

func Test_Closed(t *testing.T) {

	creditCard := NewCreditCardAccount("Credit card", 2500.00)
	bankAccount := NewAccountNoInterest("Test account", 2500.00)

	require.Equal(t, false, creditCard.Closed(), "Credit card account should not be closed")
	require.Equal(t, false, bankAccount.Closed(), "Bank account should not be closed")

	testAccount := NewBankAccount("test", TypeSavingsNoInterest, 3000.00)
	ccClosed := creditCard.Pay(testAccount)
	require.Equal(t, true, ccClosed, "Credit card account should return closed")
	require.Equal(t, true, creditCard.Closed(), "Credit card account should be closed")
	require.InDelta(t, 500.00, testAccount.Amount(), 0.01, "500.00 should be remaining in account")

	remainder := bankAccount.Deduct(4000.00)
	require.Equal(t, true, bankAccount.Closed(), "Bank account should be closed")
	require.InDelta(t, -1500.00, remainder, 0.01, "Remainder should be -1500.00")
}