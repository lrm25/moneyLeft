package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/lrm25/moneyLeft/models"
)

const FOUROHONEK_PENALTY_AGE = 59.5

const FLAG_FILE_LOCATION = "fileLocation"

const DEFAULT_FILE_LOCATION = "money.json"

const ERROR_MONEYPERMONTH = "Money per month cannot be less than or equal to 0"

type SavingsNoInterestAccount struct {
	name   string
	amount float64
}

type SavingsWithInterestAccount struct {
	name   string
	amount float64
	rate   float64
}

type StockBrokerageAccount struct {
	name         string
	amount       float64
	expectedRate float64
	expectedFee  float64
}

type FourOhOneKAccount struct {
	name              string
	amount            float64
	percentStock      float64
	expectedBondRate  float64
	expectedStockRate float64
}

type YearAndMonth struct {
	year  int
	month int
}

type Breakdown struct {
	brokeAge YearAndMonth
}

type MonthReport struct {
	time YearAndMonth
	age  YearAndMonth
	//remainingAccounts []models.Account
	milestones []Milestone
}

type Milestone struct {
	broke            bool
	year             int
	month            int
	accountEnding    string
	accountBeginning string
}

func main() {

	jsonLocation := flag.String(FLAG_FILE_LOCATION, DEFAULT_FILE_LOCATION, "JSON file containing financial data")
	flag.Parse()

	type TestStruct struct {
		TotalMoney    float64
		MoneyPerMonth float64
	}

	jsonFile, err := ioutil.ReadFile(*jsonLocation)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	var testStruct TestStruct
	if err = json.Unmarshal(jsonFile, &testStruct); err != nil {
		fmt.Printf("%v", err)
		return
	}
	if testStruct.MoneyPerMonth <= 0 {
		fmt.Printf(ERROR_MONEYPERMONTH)
		return
	}
	monthsRemaining := testStruct.TotalMoney / testStruct.MoneyPerMonth

	fmt.Printf("%.2f\n", monthsRemaining)
	dayFloatRemaining := monthsRemaining - math.Floor(monthsRemaining)
	brokeMonth := time.Now().AddDate(0, int(monthsRemaining), 0)
	fmt.Println(brokeMonth.String())
	brokeMonthOverflow := time.Date(brokeMonth.Year(), brokeMonth.Month(), 32, 0, 0, 0, 0, time.UTC)
	brokeMonthDays := 32 - brokeMonthOverflow.Day()
	brokeDay := brokeMonth.AddDate(0, 0, int(dayFloatRemaining*float64(brokeMonthDays)))
	fmt.Println(brokeDay.String())

	inflationRate := 3.28
	luke := models.NewPerson(38, 9, 90, 3619.37, inflationRate)

	stockReturnPct := 10.0
	bondReturnPct := 6.1
	midfirstBank := 18001.21
	midfirstIntRate := 1.09
	midfirstAccount := models.NewAccountWithInterest("Midfirst Bank", midfirstBank, midfirstIntRate)
	usBank := 9406.81
	usBankAccount := models.NewAccountNoInterest("US bank", usBank)
	boaInvestment := 5216.32
	boaInvestmentAccount := models.NewAccountNoInterest("Bank of America", boaInvestment)
	schwabChecking := 555.45
	schwabCheckingAccount := models.NewAccountNoInterest("Schwab Checking", schwabChecking)
	chaseSapphire := 1985.57
	chaseSapphireAccount := models.NewCreditCardAccount("Chase Sapphire", chaseSapphire)
	capitalOne := 29.00
	capitalOneAccount := models.NewCreditCardAccount("Capital One", capitalOne)
	tradIRAMoney := 108882.45
	//tradIRAMoney := 688882.45
	tradIRAPctStock := 90.0
	iraAccount := models.NewIRA("Betterment", tradIRAMoney, tradIRAPctStock, stockReturnPct, bondReturnPct, luke)
	eTrade := 24277.54
	eTradeMonthFee := 9.99
	eTradeAccount := models.NewAccountStockBrokerage("ETrade", eTrade, stockReturnPct, eTradeMonthFee, luke)
	sepIRA := 5716.89
	sepIRAPctStock := 90.0
	//sepIRAPctBond := 10
	sepIRAAccount := models.NewIRA("Betterment SEP", sepIRA, sepIRAPctStock, stockReturnPct, bondReturnPct, luke)

	schwabInvestment := 29.29
	schwabInvestmentAccount := models.NewAccountNoInterest("Schwab Investment", schwabInvestment)
	boaChecking := 578.00
	boaCheckingAccount := models.NewAccountNoInterest("Bank of America Checking", boaChecking)
	upwork := 2694.00
	upworkAccount := models.NewAccountNoInterest("Upwork", upwork)

	var creditCards []*models.CreditCardAccountImpl
	creditCards = append(creditCards, chaseSapphireAccount, capitalOneAccount)
	totalCCAmount := 0.0
	for _, card := range creditCards {
		totalCCAmount += card.Amount()
	}
	println(totalCCAmount)

	taxBrackets := models.NewTaxBrackets(13850.00,
		[]*models.TaxBracket{models.NewTaxBracket(0.0, 11000.0, 10.0),
			models.NewTaxBracket(11000.00, 44725.00, 12.0),
			models.NewTaxBracket(44725.00, 95375.00, 22.0)})
	capTaxBrackets := models.NewCapTaxBrackets(
		[]*models.TaxBracket{
			models.NewTaxBracket(0.0, 44625.00, 0.0),
			models.NewTaxBracket(44625.00, 492300.00, 15.00)})
	stateBrackets := models.NewTaxBrackets(0.00,
		[]*models.TaxBracket{models.NewTaxBracketNoMax(0.0, 4.4)})
	luke.WithTaxBrackets(taxBrackets, capTaxBrackets, stateBrackets)

	ss := models.NewPayoutsSocialSecurity(1575.00, 2261.00, 2804.00)
	ssAccount := models.NewAccountSocialSecurity(models.Early, ss, inflationRate, luke)

	ani := models.PositiveAccounts{}
	ani = append(ani, usBankAccount, boaInvestmentAccount, boaCheckingAccount, schwabCheckingAccount, schwabInvestmentAccount, upworkAccount)
	//sort.Sort(ani)

	pias := models.PassiveIncreaseAccounts{}
	pias = append(pias, midfirstAccount, eTradeAccount, sepIRAAccount, iraAccount, ssAccount)

	luke.WithAccounts(creditCards, ani, pias)
	year := 2023
	month := 8
	luke.PayCreditCards()
	if luke.Broke() {
		println("Broke from credit cards")
	}
	luke.SetIncome(1350.00)
	for {
		month++
		if month == 13 {
			month = 1
			year++
			luke.ChangeTaxYear()
		}
		println("year", year, "month", month)
		luke.IncreaseAge(year, month)
		if luke.Broke() {
			println("Broke on year ", year, " month ", month)
			return
		}
		if luke.LifeExpectancy() <= luke.AgeInYears() {
			println("You will die before you go broke")
			return
		}
	}

	/*for _, a := range ani {
		fmt.Printf("%s %.2f\n", a.GetName(), a.GetAmount())
		totalCCAmount = a.Deduct(totalCCAmount)
		if 0.0 < totalCCAmount {
			fmt.Printf("%s: %.2f remaining\n", a.GetName(), a.GetAmount())
			break
		} else {
			println("Used up " + a.GetName())
			ani = ani[1:]
			totalCCAmount = math.Abs(totalCCAmount)
		}
	}
	if len(ani) == 0 {
		println("you already have a negative net worth")
		return
	}

	currentAgeYears := 38
	currentAgeMonths := 8
	currentYear := 2023
	currentMonth := 7
	lifeExpectancy := 90
	//totalMoney := testStruct.TotalMoney
	//totalMoney -= fourOhOneKMoney
	moneyPerMonth := testStruct.MoneyPerMonth
	var moneyNeeded float64
	programComplete := false
	for !programComplete {
		monthComplete := false
		for !monthComplete {
			fmt.Printf("Year %d month %d", currentYear, currentMonth)
			moneyNeeded = moneyPerMonth
			var moneyRemaining float64
			for _, a := range ani {
				moneyRemaining = a.Deduct(moneyNeeded)
				fmt.Printf("account: %s, money remaining: %.2f\n", a.GetName(), moneyRemaining)
				if 0.0 < moneyRemaining {
					monthComplete = true
					break
				} else {
					moneyNeeded = math.Abs(moneyRemaining)
					ani = ani[1:]
					if len(ani) == 0 {
						fmt.Printf("Run out of money in year %d month %d, age %d and %d months\n", currentYear, currentMonth, currentAgeYears, currentAgeMonths)
						return
					}
				}
			}
		}
		currentAgeMonths++
		if currentAgeMonths == 12 {
			currentAgeYears++
			currentAgeMonths = 0
		}
		currentMonth++
		if 12 < currentMonth {
			currentYear++
			currentMonth = 1
		}
		if lifeExpectancy <= currentAgeYears {
			println("You won't run out of money")
			break
		}
		for _, pia := range pias {
			if !pia.Closed() {
				println("before increase", pia.GetAmount())
				pia.Increase()
				println("after increase", pia.GetAmount())
			}
		}
		moneyPerMonth *= (1 + inflationRate/1200)
		//fmt.Printf("remaining: %d %d %d %f\n", currentMonth, currentAgeYears, currentAgeMonths, totalMoney)
	}*/

	/*fourOhOneKMoney -= moneyNeeded

	tax := 0.12
	for 0 <= fourOhOneKMoney {
		penalty := 0.0
		if currentAgeYears < 60 || (currentAgeYears == 59 && currentAgeMonths < 6) {
			penalty = 0.1
		}
		fourOhOneKMoney -= (moneyPerMonth / (1 - (tax + penalty)))
		if fourOhOneKMoney < 0 {
			break
		}
		moneyPerMonth *= (1 + 0.04/12)
		fourOhOneKMoney *= (1 + 0.1/12)
		currentAgeMonths++
		if currentAgeMonths == 12 {
			currentAgeYears++
			currentAgeMonths = 0
		}
		currentMonth++
		if 12 < currentMonth {
			currentMonth = 1
		}
		if lifeExpectancy <= currentAgeYears {
			break
		}
		fmt.Printf("401k: %d %d %d %f\n", currentMonth, currentAgeYears, currentAgeMonths, fourOhOneKMoney)
	}*/
	//println("Broke age", currentAgeYears, currentAgeMonths)

	/*var inflationRate float64
	var stockReturnRate float64
	var bondReturnRate float64
	var socialSecurity SocialSecurityAccount
	var leisureMultiplier float64
	var moneyNeededPerMonth float64*/

	// load json file
	// import data from json file
	// return error if failure
	// return results
	// overall time remaining
	// month-by-month, with expenses
}
