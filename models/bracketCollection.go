package models

// TaxBracketCollection is a single struct containing all relevant tax brackets
type TaxBracketCollection struct {
	federal      *FedTaxBrackets
	state        *StateTaxBrackets
	capitalGains *CapTaxBrackets
}

// NewTaxBracketCollection constructs a tax bracket collection from the different types of tax brackets
func NewTaxBracketCollection(federal *FedTaxBrackets, state *StateTaxBrackets, capitalGains *CapTaxBrackets) *TaxBracketCollection {
	return &TaxBracketCollection{
		federal:      federal,
		state:        state,
		capitalGains: capitalGains,
	}
}

// Federal gets US federal tax brackets
func (c *TaxBracketCollection) Federal() *FedTaxBrackets {
	return c.federal
}

// State gets US state tax brackets
func (c *TaxBracketCollection) State() *StateTaxBrackets {
	return c.state
}

// CapitalGains gets US capital gains tax brackets
func (c *TaxBracketCollection) CapitalGains() *CapTaxBrackets {
	return c.capitalGains
}
