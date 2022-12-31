package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/lrm25/moneyLeft/account"
)

const FOUROHONEK_PENALTY_AGE = 59.5

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
	time              YearAndMonth
	age               YearAndMonth
	remainingAccounts []account.Account
	milestones        []Milestone
}

type Milestone struct {
	broke            bool
	year             int
	month            int
	accountEnding    string
	accountBeginning string
}

func main() {

	type TestStruct struct {
		TotalMoney    float64
		MoneyPerMonth float64
	}

	jsonFile, err := ioutil.ReadFile("money.json")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	var testStruct TestStruct
	if err = json.Unmarshal(jsonFile, &testStruct); err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%.2f\n", testStruct.TotalMoney/testStruct.MoneyPerMonth)

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
