package models

// TaxBracket is a struct representing a US tax bracket
type TaxBracket struct {
	minimum    float64
	hasMaximum bool
	maximum    float64
	rate       float64
}

// NewTaxBracket is a constructor for tax brackets with a maximum amount
func NewTaxBracket(minimum, maximum, rate float64) *TaxBracket {
	return &TaxBracket{
		minimum:    minimum,
		maximum:    maximum,
		rate:       rate,
		hasMaximum: true,
	}
}

// NewTaxBracketNoMax is a constructor for tax brackets without a maximum amount
func NewTaxBracketNoMax(minimum, rate float64) *TaxBracket {
	return &TaxBracket{
		minimum:    minimum,
		hasMaximum: false,
		rate:       rate,
	}
}

// HasMaximum returns true if this bracket does not have a maximum amount before transition to a higher bracket.
func (t *TaxBracket) HasMaximum() bool {
	return t.hasMaximum
}

// Minimum returns the minimum income that this tax bracket covers.
func (t *TaxBracket) Minimum() float64 {
	return t.minimum
}

// Maximum returns the maximum income that this tax bracket covers.  If HasMaximum is set, this should be ignored.
func (t *TaxBracket) Maximum() float64 {
	return t.maximum
}

// Rate returns the tax rate owed for this bracket as a decimal below 1
func (t *TaxBracket) Rate() float64 {
	return t.rate
}

type brackets []*TaxBracket

// Brackets returns the tax brackets held by this struct
func (b *brackets) Brackets() []*TaxBracket {
	return *b
}

// Inflate adjusts brackets based on expected inflation rate
func (b *brackets) Inflate(inflationRate float64) {
	for _, bracket := range *b {
		bracket.minimum *= (1 + inflationRate/100)
		if bracket.hasMaximum {
			bracket.maximum *= (1 + inflationRate/100)
		}
	}
}