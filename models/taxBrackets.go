package models

type TaxBracket struct {
	Minimum    float64
	HasMaximum bool
	Maximum    float64
	Rate       float64
}

func NewTaxBracket(minimum, maximum, rate float64) *TaxBracket {
	return &TaxBracket{
		Minimum:    minimum,
		Maximum:    maximum,
		Rate:       rate,
		HasMaximum: true,
	}
}

func NewTaxBracketNoMax(minimum, rate float64) *TaxBracket {
	return &TaxBracket{
		Minimum:    minimum,
		HasMaximum: false,
		Rate:       rate,
	}
}

type TaxBrackets struct {
	StandardDeduction float64
	Brackets          []*TaxBracket
}

type CapTaxBrackets struct {
	Brackets []*TaxBracket
}

func NewCapTaxBrackets(brackets []*TaxBracket) *CapTaxBrackets {
	return &CapTaxBrackets{
		Brackets: brackets,
	}
}

func NewTaxBrackets(standardDeduction float64, brackets []*TaxBracket) *TaxBrackets {
	return &TaxBrackets{
		StandardDeduction: standardDeduction,
		Brackets:          brackets,
	}
}

func (t *TaxBrackets) GetTaxAmount(yearlyIncome float64) float64 {
	taxAmount := 0.00
	moneyRemaining := yearlyIncome - t.StandardDeduction
	if moneyRemaining <= 0 {
		return taxAmount
	}
	for _, bracket := range t.Brackets {
		if !bracket.HasMaximum || moneyRemaining < bracket.Maximum {
			taxAmount += (moneyRemaining - bracket.Minimum) * bracket.Rate / 100
			return taxAmount
		}
		taxAmount += bracket.Rate / 100 * (bracket.Maximum - bracket.Minimum)
	}
	return taxAmount
}

func (t *CapTaxBrackets) GetTaxAmount(nonCapTaxIncome, capTaxIncome float64) float64 {
	taxAmount := 0.00
	if capTaxIncome < 0.01 {
		return taxAmount
	}
	for _, bracket := range t.Brackets {
		adjustedBracketMin := bracket.Minimum - nonCapTaxIncome
		adjustedBracketMax := 0.00
		if bracket.HasMaximum {
			adjustedBracketMax = bracket.Maximum - nonCapTaxIncome
			if adjustedBracketMax <= 0 {
				continue
			}
		}
		if adjustedBracketMin < 0 {
			adjustedBracketMin = 0
		}
		if !bracket.HasMaximum || capTaxIncome < adjustedBracketMax {
			taxAmount += (capTaxIncome - adjustedBracketMin) * bracket.Rate / 100
			println(taxAmount)
			return taxAmount
		}
		taxAmount += bracket.Rate / 100 * (adjustedBracketMax - adjustedBracketMin)
		println(taxAmount)
	}
	return taxAmount
}

func (t *TaxBrackets) Inflate(inflationRate float64) {
	t.StandardDeduction *= (1 + inflationRate/100)
	for _, b := range t.Brackets {
		b.Minimum *= (1 + inflationRate/100)
		if b.HasMaximum {
			b.Maximum *= (1 + inflationRate/100)
		}
	}
}

func (t *CapTaxBrackets) Inflate(inflationRate float64) {
	println(t.Brackets)
	for _, b := range t.Brackets {
		b.Minimum *= (1 + inflationRate/100)
		if b.HasMaximum {
			b.Maximum *= (1 + inflationRate/100)
		}
	}
}
