package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

// CapTaxBrackets struct
type CapTaxBrackets struct {
	brackets
}

// NewCapTaxBrackets returns capital tax struct from provided brackets
func NewCapTaxBrackets(brackets []*TaxBracket) *CapTaxBrackets {
	return &CapTaxBrackets{
		brackets: brackets,
	}
}

// GetTaxAmount gets yearly tax amount based on non-capital gains and capital gains income, excluding the former from taxation,
// but factoring it into the tax rate
func (t *CapTaxBrackets) GetTaxAmount(nonCapTaxIncome, capTaxIncome float64) float64 {

	taxAmount := 0.00
	if capTaxIncome < 0.01 {
		return taxAmount
	}
	for _, bracket := range t.brackets {
		adjustedBracketMin := bracket.minimum - nonCapTaxIncome
		adjustedBracketMax := 0.00
		if bracket.hasMaximum {
			adjustedBracketMax = bracket.maximum - nonCapTaxIncome
			if adjustedBracketMax <= 0 {
				continue
			}
		}
		if adjustedBracketMin < 0 {
			adjustedBracketMin = 0
		}
		if !bracket.hasMaximum || capTaxIncome < adjustedBracketMax {
			taxAmount += (capTaxIncome - adjustedBracketMin) * bracket.rate / 100
			logger.Get().Debug(fmt.Sprintf("Cap tax amount: %.2f", taxAmount))
			return taxAmount
		}
		taxAmount += bracket.rate / 100 * (adjustedBracketMax - adjustedBracketMin)
		logger.Get().Debug(fmt.Sprintf("Cap tax amount: %.2f", taxAmount))
	}
	return taxAmount
}