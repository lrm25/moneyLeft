package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TaxBracketsStdDeduction(t *testing.T) {

	standardDeduction := 10000.00
	taxBracket := NewTaxBracketNoMax(0, 10.00)
	taxBrackets := NewTaxBrackets(standardDeduction, []*TaxBracket{taxBracket})

	yearlyIncome := 11000.00
	require.InDelta(t, 100.00, taxBrackets.GetTaxAmount(yearlyIncome), 1.00, "Tax obligation should be close to 100.00")
}

func Test_TaxBracketsCapGains(t *testing.T) {

	taxBracket := NewTaxBracket(0.00, 20000.00, 10.0)
	taxBracketTwo := NewTaxBracketNoMax(20000.00, 20.0)
	taxBrackets := NewCapTaxBrackets([]*TaxBracket{taxBracket, taxBracketTwo})

	nonCapGainsAmount := 10000.00
	capGainsAmount := 20000.00
	amount := taxBrackets.GetTaxAmount(nonCapGainsAmount, capGainsAmount)
	require.InDelta(t, 3000.00, amount, 1.00)
}
