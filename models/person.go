package models

import (
	"fmt"

	"github.com/lrm25/moneyLeft/logger"
)

// Person - struct representing person running program, and info
type Person struct {
	years               int
	months              int
	lifeExpectancy      int
	neededPerMonth      float64
	taxableCapGainsLast float64
	taxableOtherLast    float64
	taxableCapGainsThis float64
	taxableOtherThis    float64
	inflationRate       float64
	broke               bool
	alive               bool
	creditCards         []*CreditCardAccountImpl
	accounts            []PositiveAccount
	interestAccounts    []PassiveIncreaseAccount
	nonCapBrackets      *FedTaxBrackets
	capBrackets         *CapTaxBrackets
	stateBrackets       *StateTaxBrackets
	income              float64
}

// NewPerson construtor (current age in years and months, life expectancy, money needed per month, expected yearly inflation rate)
func NewPerson(years, months, lifeExpectancy int, neededPerMonth, inflationRate float64) *Person {
	return &Person{
		years:          years,
		months:         months,
		lifeExpectancy: lifeExpectancy,
		neededPerMonth: neededPerMonth,
		broke:          false,
		alive:          true,
		inflationRate:  inflationRate,
	}
}

// WithAccounts sets the person's accounts (credit cards, non-passively-increasing accounts, passively-increasing accounts)
func (p *Person) WithAccounts(creditCards []*CreditCardAccountImpl, accounts PositiveAccounts, interestAccounts PassiveIncreaseAccounts) *Person {
	p.creditCards = creditCards
	p.accounts = accounts
	p.interestAccounts = interestAccounts
	return p
}

// WithTaxBrackets sets the person's expected tax brackets
func (p *Person) WithTaxBrackets(nonCapBrackets *FedTaxBrackets, capBrackets *CapTaxBrackets, stateBrackets *StateTaxBrackets) {
	p.nonCapBrackets = nonCapBrackets
	p.capBrackets = capBrackets
	p.stateBrackets = stateBrackets
}

// Broke returns whether or not the user is broke.  If so, the program should terminate.
func (p *Person) Broke() bool {
	return p.broke
}

// LifeExpectancy returns the user's life expectancy.  If the user reaches their life expectancy before going broke, this
// will be reported by the program.
func (p *Person) LifeExpectancy() int {
	return p.lifeExpectancy
}

// AgeYears returns the user's age in years as this program runs
func (p *Person) AgeYears() int {
	return p.years
}

// AgeMonths returns the age in months following the age in years
func (p *Person) AgeMonths() int {
	return p.months
}

// NeededPerMonth returns the user's monthly expense amount
func (p *Person) NeededPerMonth() float64 {
	return p.neededPerMonth
}

// SetIncome sets a user's monthly income
func (p *Person) SetIncome(income float64) {
	p.income = income
}

// PayCreditCards immediately pays off the user's credit cards.  To keep things simple now,
// the user is immediately declared broke if they owe more than they have.
func (p *Person) PayCreditCards() {
	if len(p.creditCards) == 0 {
		return
	}
	if len(p.accounts) == 0 {
		p.broke = true
		return
	}
	for _, creditCard := range p.creditCards {
		for idx, account := range p.accounts {
			if creditCard.Pay(account) {
				break
			}
			if idx == len(p.accounts)-1 {
				break
			}
		}
		for idx, account := range p.interestAccounts {
			if creditCard.Pay(account) {
				break
			}
			if idx == len(p.accounts)-1 {
				break
			}
		}
		if 0 < creditCard.Amount() {
			p.broke = true
			return
		}
	}
	p.creditCards = nil
}

// pay the user's monthly expenses
func (p *Person) pay(remaining float64) bool {
	for _, account := range p.accounts {
		if !account.Closed() {
			logger.Get().Debug(fmt.Sprintf("Paying from account %s", account.Name()))
			remaining = account.Deduct(remaining)
			logger.Get().Debug(fmt.Sprintf("non interest remaining: %.2f", remaining))
			if 0 < remaining {
				return true
			}
			remaining *= -1
		}
	}

	for _, account := range p.interestAccounts {
		if !account.Closed() {
			logger.Get().Debug(fmt.Sprintf("Paying from account %s", account.Name()))
			logger.Get().Debug(fmt.Sprintf("remaining: %.2f", remaining))
			remaining = account.Deduct(remaining)
			logger.Get().Debug(fmt.Sprintf("after deduction: %.2f", remaining))
			if 0 < remaining {
				return true
			}
			remaining *= -1
		}
	}
	return false
}

// ChangeTaxYear switches the tax year and begins incrementing taxes on a new year, leaving
// the previous year's taxable income to be available to be deducted on the tax month.
func (p *Person) ChangeTaxYear() {
	p.taxableCapGainsLast = p.taxableCapGainsThis
	p.taxableOtherLast = p.taxableOtherThis
	p.taxableCapGainsThis = 0
	p.taxableOtherThis = 0
}

// PayTaxes deducts the user's expected yearly taxes
func (p *Person) PayTaxes() bool {

	amount := 0.00
	if p.nonCapBrackets != nil {
		amount += p.nonCapBrackets.GetTaxAmount(p.taxableOtherLast)
		logger.Get().Debug(fmt.Sprintf("non cap amount: %.2f", p.nonCapBrackets.GetTaxAmount(p.taxableOtherLast)))
	}
	if p.capBrackets != nil {
		amount += p.capBrackets.GetTaxAmount(p.taxableOtherLast, p.taxableCapGainsLast)
		logger.Get().Debug(fmt.Sprintf("cap amount: %.2f", p.capBrackets.GetTaxAmount(p.taxableOtherLast, p.taxableCapGainsLast)))
	}
	if p.stateBrackets != nil {
		amount += p.stateBrackets.GetTaxAmount(p.taxableOtherLast+p.taxableCapGainsLast)
		logger.Get().Debug(fmt.Sprintf("state amount: %.2f", p.stateBrackets.GetTaxAmount(p.taxableOtherLast+p.taxableCapGainsLast)))
	}

	logger.Get().Debug(fmt.Sprintf("amount: %.2f", amount))
	// reset afterwards
	p.taxableCapGainsLast = 0
	p.taxableOtherLast = 0
	return p.pay(amount)
}

// IncreaseAge increases the person's age by one month and does corresponding calculations:  deducting
// expected monthly amount, closing empty accounts with the exception of social security, paying taxes
// once a year, and terminating if the user is broke.  The program also allows the user to add a monthly
// income if they make less than what they need and want to see how much longer they can last if they 
// continue to make that income.
func (p *Person) IncreaseAge(year, month int) {
	p.months++
	if p.months == 12 {
		p.years++
		p.months = 0
	}
	p.neededPerMonth *= (1 + p.inflationRate/1200.00)
	logger.Get().Debug(fmt.Sprintf("%.2f", p.neededPerMonth))

	neededPerMonth := p.neededPerMonth
	// for now, keep this simple and don't worry about handling first year YTD taxes
	if 0 < p.income {
		neededPerMonth -= p.income
		p.taxableOtherThis += p.income
		p.income *= (1 + p.inflationRate/1200.00)
	}

	for _, account := range p.interestAccounts {
		if !account.Closed() || account.Type() == TypeSocialSecurity {
			account.Increase()
		}
	}
	if !p.pay(neededPerMonth) {
		p.broke = true
		return
	}
	for idx:= 0; idx < len(p.accounts); {
		if p.accounts[idx].Closed() && p.accounts[idx].Removable() {
			logger.Get().Debug(fmt.Sprintf("Removing account %s", p.accounts[idx].Name()))
			if idx < len(p.accounts)-1 {
				p.accounts = append(p.accounts[:idx], p.accounts[idx+1:]...)
			} else {
				p.accounts = p.accounts[:idx]
			}
			continue
		}
		idx++
	}
	for idx := 0; idx < len(p.accounts); {
		if p.accounts[idx].Closed() && p.accounts[idx].Removable() {
			logger.Get().Debug(fmt.Sprintf("Removing account %s", p.accounts[idx].Name()))
			if idx < len(p.interestAccounts)-1 {
				p.interestAccounts = append(p.interestAccounts[:idx], p.interestAccounts[idx+1:]...)
			} else {
				p.interestAccounts = p.interestAccounts[:idx]
			}
			continue
		}
		idx++
	}
	if month == 4 {
		if !p.PayTaxes() {
			p.broke = true
			return
		}
		p.nonCapBrackets.Inflate(p.inflationRate)
		p.capBrackets.Inflate(p.inflationRate)
		p.stateBrackets.Inflate(p.inflationRate)
	}
	logger.Get().Debug(p.String(year, month))
}

// String returns a string representation of the person for info, debugging
func (p *Person) String(year, month int) string {
	s := fmt.Sprintf("STATUS - year %d month %d\n", year, month)
	s += fmt.Sprintf("Age: %d years %d months\n", p.years, p.months)
	s += fmt.Sprintf("Needed per month: %.2f", p.neededPerMonth)
	for _, cc := range p.creditCards {
		s += fmt.Sprintf("\n%s %.2f", cc.name, cc.amount)
	}
	for _, account := range p.accounts {
		s += fmt.Sprintf("\n%s", account.String())
	}
	for _, pa := range p.interestAccounts {
		s += fmt.Sprintf("\n%s", pa.String())
	}
	return s
}