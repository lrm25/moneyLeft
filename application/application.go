package application

import (
	"fmt"
	"time"

	"github.com/lrm25/moneyLeft/config"
	"github.com/lrm25/moneyLeft/logger"
	"github.com/lrm25/moneyLeft/models"
)

func runAgeLoop(person *models.Person, year, month int) {
	totalMonths := 0
	for {
		month++
		totalMonths++
		if month == 13 {
			month = 1
			year++
			person.ChangeTaxYear()
		}
		logger.Get().Debug(fmt.Sprintf("year: %d, month: %d", year, month))
		person.IncreaseAge(year, month)
		if person.Broke() {
			logger.Get().Info(fmt.Sprintf("Broke on year %d, month %d", year, month))
			logger.Get().Info(fmt.Sprintf("Age:  %d years, %d months", person.AgeYears(), person.AgeMonths()))
			logger.Get().Info(fmt.Sprintf("Total time:  %d years, %d months", totalMonths/12, totalMonths%12))
			return
		}
		if person.LifeExpectancy() <= person.AgeYears() {
			logger.Get().Info("You will die before you go broke")
			return
		}
	}
}

// Run the application, telling the user how long they can last before going broke with current accounts and monthly expenses
func Run(c *config.YamlConfig) {

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
	runAgeLoop(person, year, month)
}
