package models

// StateTaxBrackets struct
type StateTaxBrackets struct {
	brackets
}

// NewStateTaxBrackets constructor
func NewStateTaxBrackets(brackets []*TaxBracket) *StateTaxBrackets {
	return &StateTaxBrackets{
		brackets: brackets,
	}
}

// GetTaxAmount returns the yearly state tax owed based on total income (non-cap gains + cap gains)
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