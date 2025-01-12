package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AccountStockBrokerage_noPerson(t *testing.T) {
	require.Panics(t, func() { NewAccountStockBrokerage("name", 0.0, 0.0, 0.0, nil) })
}

func Test_AccountStockBrokerage_deduct(t *testing.T) {
	person := NewPerson(50, 0, 60, 1000.00, 5.00)
	asb := NewAccountStockBrokerage("name", 40.00, 0.0, 10.0, person)
	asb.Deduct(10.00)
	require.InDelta(t, 20.0, asb.amount, 0.01)
	require.False(t, asb.Closed())
	asb.Deduct(10.01)
	require.True(t, asb.Closed())
}
