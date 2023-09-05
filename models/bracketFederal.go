package models

// FedTaxBrackets contains the US standard deduction and tax brackets
type FedTaxBrackets struct {
	brackets
	standardDeduction float64
}

// NewFedTaxBrackets returns a federal tax struct from the specified standard deduction and brackets
func NewFedTaxBrackets(standardDeduction float64, brackets []*TaxBracket) *FedTaxBrackets {
	return &FedTaxBrackets{
		standardDeduction: standardDeduction,
		brackets:          brackets,
	}
}

// StandardDeduction returns US federal standard deduction
func (t *FedTaxBrackets) StandardDeduction() float64 {
	return t.standardDeduction
}

// GetTaxAmount retrieves yearly tax amount based on taxable income and federal tax brackets
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