package application

import (
	"fmt"
	"time"

	"github.com/lrm25/moneyLeft/config"
	"github.com/lrm25/moneyLeft/logger"
	"github.com/lrm25/moneyLeft/models"
)

type applicationInfo struct {
	person             *models.Person
	year               int
	month              int
	totalMonths        int
	calculateMinNeeded bool
}

func (a *applicationInfo) runSingleMonth() (bool, bool) {
	a.month++
	a.totalMonths++
	if a.month == 13 {
		a.month = 1
		a.year++
		a.person.ChangeTaxYear()
	}
	logger.Get().Debug(fmt.Sprintf("year: %d, month: %d", a.year, a.month))
	a.person.IncreaseAge(a.year, a.month)
	if a.person.Broke() {
		if !a.calculateMinNeeded {
			logger.Get().Info(fmt.Sprintf("Broke on year %d, month %d", a.year, a.month))
			logger.Get().Info(fmt.Sprintf("Age:  %d years, %d months", a.person.AgeYears(), a.person.AgeMonths()))
			logger.Get().Info(fmt.Sprintf("Total time:  %d years, %d months", a.totalMonths/12, a.totalMonths%12))
		}
		return true, true
	}
	if a.person.LifeExpectancy() <= a.person.AgeYears() {
		logger.Get().Info("You will die before you go broke")
		return false, true
	}
	return false, false
}

func runAgeLoop(person *models.Person, year, month int, calculateMinNeeded bool) bool {
	appInfo := &applicationInfo{
		person:             person,
		year:               year,
		month:              month,
		calculateMinNeeded: calculateMinNeeded,
		totalMonths:        0,
	}
	broke := false
	complete := false
	for {
		broke, complete = appInfo.runSingleMonth()
		if complete {
			break
		}
	}
	return broke
}

// Run the application, telling the user how long they can last before going broke with current accounts and monthly expenses
func Run(c *config.YamlConfig, calculateMinNeeded bool) bool {

	person := c.Person()

	person.WithTaxBrackets(c.TaxBracketCollection().Federal(), c.TaxBracketCollection().CapitalGains(), c.TaxBracketCollection().State())
	pas := models.PositiveAccounts{}
	for _, a := range c.NoInterestAccounts() {
		pas = append(pas, a)
	}

	pias := models.PassiveIncreaseAccounts{}
	for _, a := range c.InterestAccounts() {
		pias = append(pias, a)
	}
	for _, a := range c.RealEstateInvestments() {
		pias = append(pias, a)
	}
	for _, a := range c.BrokerageAccounts() {
		pias = append(pias, a)
	}
	for _, a := range c.IRAs() {
		pias = append(pias, a)
	}
	pias = append(pias, c.SocialSecurity())

	person.WithAccounts(c.CreditCards(), pas, pias)
	year := time.Now().Year()
	month := int(time.Now().Month())

	if c.MonthlyIncome != nil {
		person.SetIncome(*c.MonthlyIncome)
	}

	person.PayCreditCards()
	if person.Broke() {
		logger.Get().Info("Broke from credit cards")
	}
	return runAgeLoop(person, year, month, calculateMinNeeded)
}
