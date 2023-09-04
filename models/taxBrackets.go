package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

type TaxBracket struct {
	minimum    float64
	hasMaximum bool
	maximum    float64
	rate       float64
}

func NewTaxBracket(minimum, maximum, rate float64) *TaxBracket {
	return &TaxBracket{
		minimum:    minimum,
		maximum:    maximum,
		rate:       rate,
		hasMaximum: true,
	}
}

func NewTaxBracketNoMax(minimum, rate float64) *TaxBracket {
	return &TaxBracket{
		minimum:    minimum,
		hasMaximum: false,
		rate:       rate,
	}
}

func (t *TaxBracket) HasMaximum() bool {
	return t.hasMaximum
}

func (t *TaxBracket) Minimum() float64 {
	return t.minimum
}

func (t *TaxBracket) Maximum() float64 {
	return t.maximum
}

func (t *TaxBracket) Rate() float64 {
	return t.rate
}

type brackets []*TaxBracket

type FedTaxBrackets struct {
	brackets
	standardDeduction float64
}

type CapTaxBrackets struct {
	brackets
}

type StateTaxBrackets struct {
	brackets
}

type TaxBracketCollection struct {
	federal      *FedTaxBrackets
	state        *StateTaxBrackets
	capitalGains *CapTaxBrackets
}

func NewTaxBracketCollection(federal *FedTaxBrackets, state *StateTaxBrackets, capitalGains *CapTaxBrackets) *TaxBracketCollection {
	return &TaxBracketCollection{
		federal:      federal,
		state:        state,
		capitalGains: capitalGains,
	}
}

func NewFedTaxBrackets(standardDeduction float64, brackets []*TaxBracket) *FedTaxBrackets {
	return &FedTaxBrackets{
		standardDeduction: standardDeduction,
		brackets:          brackets,
	}
}

func NewStateTaxBrackets(brackets []*TaxBracket) *StateTaxBrackets {
	return &StateTaxBrackets{
		brackets: brackets,
	}
}

func NewCapTaxBrackets(brackets []*TaxBracket) *CapTaxBrackets {
	return &CapTaxBrackets{
		brackets: brackets,
	}
}

func (b *brackets) Brackets() []*TaxBracket {
	return *b
}

func (c *TaxBracketCollection) Federal() *FedTaxBrackets {
	return c.federal
}

func (c *TaxBracketCollection) State() *StateTaxBrackets {
	return c.state
}

func (c *TaxBracketCollection) CapitalGains() *CapTaxBrackets {
	return c.capitalGains
}

func (t *FedTaxBrackets) StandardDeduction() float64 {
	return t.standardDeduction
}

func (t *FedTaxBrackets) GetTaxAmount(yearlyIncome float64) float64 {
	taxAmount := 0.00
	moneyRemaining := yearlyIncome - t.standardDeduction
	if moneyRemaining <= 0 {
		return taxAmount
	}
	for _, bracket := range t.brackets {
		if !bracket.hasMaximum || moneyRemaining < bracket.maximum {
			taxAmount += (moneyRemaining - bracket.minimum) * bracket.rate / 100
			return taxAmount
		}
		taxAmount += bracket.rate / 100 * (bracket.maximum - bracket.minimum)
	}
	return taxAmount
}

func (t *StateTaxBrackets) GetTaxAmount(total float64) float64 {
	taxAmount := 0.00
	for _, bracket := range t.brackets {
		if !bracket.hasMaximum || total < bracket.maximum {
			taxAmount += (total - bracket.minimum) * bracket.rate / 100
			return taxAmount
		}
		taxAmount += bracket.rate / 100 * (bracket.maximum - bracket.minimum)
	}
	return taxAmount
}

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

/*func (t *FedTaxBrackets) Inflate(inflationRate float64) {
	t.standardDeduction *= (1 + inflationRate/100)
	for _, b := range t.brackets {
		b.minimum *= (1 + inflationRate/100)
		if b.hasMaximum {
			b.maximum *= (1 + inflationRate/100)
		}
	}
}*/

func (t *brackets) Inflate(inflationRate float64) {
	for _, b := range *t {
		b.minimum *= (1 + inflationRate/100)
		if b.hasMaximum {
			b.maximum *= (1 + inflationRate/100)
		}
	}
}