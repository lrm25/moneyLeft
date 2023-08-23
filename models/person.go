package models

import "fmt"

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
	nonCapBrackets      *TaxBrackets
	capBrackets         *CapTaxBrackets
	stateBrackets       *TaxBrackets
	income              float64
}

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

func (p *Person) WithAccounts(creditCards []*CreditCardAccountImpl, accounts PositiveAccounts, interestAccounts PassiveIncreaseAccounts) *Person {
	p.creditCards = creditCards
	p.accounts = accounts
	p.interestAccounts = interestAccounts
	return p
}

func (p *Person) WithTaxBrackets(nonCapBrackets *TaxBrackets, capBrackets *CapTaxBrackets, stateBrackets *TaxBrackets) {
	p.nonCapBrackets = nonCapBrackets
	p.capBrackets = capBrackets
	p.stateBrackets = stateBrackets
}

func (p *Person) Broke() bool {
	return p.broke
}

func (p *Person) LifeExpectancy() int {
	return p.lifeExpectancy
}

func (p *Person) AgeInYears() int {
	return p.years
}

func (p *Person) SetIncome(income float64) {
	p.income = income
}

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

func (p *Person) pay(remaining float64) bool {
	for _, account := range p.accounts {
		if !account.Closed() {
			fmt.Printf("Paying from account %s\n", account.Name())
			remaining = account.Deduct(remaining)
			fmt.Printf("non interest remaining: %.2f\n", remaining)
			if 0 < remaining {
				return true
			} else {
				remaining *= -1
			}
		}
	}

	for _, account := range p.interestAccounts {
		if !account.Closed() {
			fmt.Printf("Paying from account %s\n", account.Name())
			fmt.Printf("remaining: %.2f\n", remaining)
			remaining = account.Deduct(remaining)
			fmt.Printf("after deduction: %.2f\n", remaining)
			if 0 < remaining {
				return true
			} else {
				remaining *= -1
			}
		}
	}
	return false
}

func (p *Person) ChangeTaxYear() {
	p.taxableCapGainsLast = p.taxableCapGainsThis
	p.taxableOtherLast = p.taxableOtherThis
	p.taxableCapGainsThis = 0
	p.taxableOtherThis = 0
}

func (p *Person) PayTaxes() bool {
	amount := p.nonCapBrackets.GetTaxAmount(p.taxableOtherLast) + p.capBrackets.GetTaxAmount(p.taxableOtherLast, p.taxableCapGainsLast) + p.stateBrackets.GetTaxAmount(p.taxableOtherLast+p.taxableCapGainsLast)
	fmt.Printf("non cap amount: %.2f\n", p.nonCapBrackets.GetTaxAmount(p.taxableOtherLast))
	fmt.Printf("cap amount: %.2f\n", p.capBrackets.GetTaxAmount(p.taxableOtherLast, p.taxableCapGainsLast))
	fmt.Printf("state amount: %.2f\n", p.stateBrackets.GetTaxAmount(p.taxableOtherLast+p.taxableCapGainsLast))
	fmt.Printf("amount: %.2f\n", amount)
	// reset afterwards
	p.taxableCapGainsLast = 0
	p.taxableOtherLast = 0
	return p.pay(amount)
}

func (p *Person) IncreaseAge(year, month int) {
	p.months++
	if p.months == 13 {
		p.years++
		p.months = 1
	}
	p.neededPerMonth *= (1 + p.inflationRate/1200.00)
	println(p.neededPerMonth)

	neededPerMonth := p.neededPerMonth
	if 0 < p.income {
		neededPerMonth -= p.income
		p.taxableOtherThis += p.income
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
	for idx, account := range p.accounts {
		if account.Closed() && account.Removable() {
			fmt.Printf("Removing account %s\n", account.Name())
			if idx < len(p.accounts)-1 {
				p.accounts = append(p.accounts[:idx], p.accounts[idx+1:]...)
			} else {
				p.accounts = p.accounts[:idx]
			}
			idx--
		}
	}
	for idx, account := range p.interestAccounts {
		if account.Closed() && account.Removable() {
			fmt.Printf("Removing account %s\n", account.Name())
			if idx < len(p.interestAccounts)-1 {
				p.interestAccounts = append(p.interestAccounts[:idx], p.interestAccounts[idx+1:]...)
			} else {
				p.interestAccounts = p.interestAccounts[:idx]
			}
			idx--
		}
	}
	if month == 4 {
		if !p.PayTaxes() {
			p.broke = true
			return
		}
		p.nonCapBrackets.Inflate(p.inflationRate)
		p.capBrackets.Inflate(p.inflationRate)
	}
	p.PrintStatus(year, month)
}

func (p *Person) PrintStatus(year, month int) {
	fmt.Printf("STATUS - year %d month %d\n", year, month)
	fmt.Printf("Age: %d years %d months\n", p.years, p.months)
	fmt.Printf("Needed per month: %.2f\n", p.neededPerMonth)
	for _, cc := range p.creditCards {
		fmt.Printf("%s %.2f\n", cc.name, cc.amount)
	}
	for _, account := range p.accounts {
		fmt.Printf("%s\n", account.String())
	}
	for _, pa := range p.interestAccounts {
		fmt.Printf("%s\n", pa.String())
	}
}
