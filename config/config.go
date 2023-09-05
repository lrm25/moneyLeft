package config

import (
	"fmt"
	"os"

	"github.com/lrm25/moneyLeft/logger"
	"github.com/lrm25/moneyLeft/models"
	"gopkg.in/yaml.v2"
)

// Constants representing social security payout declaration in yaml file
const (
	SSPayoutEarly = "early"
	SSPayoutNormal = "normal"
	SSPayoutLate = "late"
)

// YamlPerson - person in yaml file
type YamlPerson struct {
	Years          int `yaml:"years"`
	Months         int `yaml:"months"`
	LifeExpectancy int `yaml:"lifeExpectancy"`
	NeededPerMonth float64 `yaml:"neededPerMonth"`
}

// CreditCard - credit card in yaml file
type CreditCard struct {
	Name   string  `yaml:"name"`
	Amount float64 `yaml:"amount"`
}

// InterestAccount - passive increase account in yaml file
type InterestAccount struct {
	Name   string  `yaml:"name"`
	Amount float64 `yaml:"amount"`
	Rate   float64 `yaml:"rate"`
}

// NoInterestAccount - no passive increase account in yaml file
type NoInterestAccount struct {
	Name   string  `yaml:"name"`
	Amount float64 `yaml:"amount"`
}

// BrokerageAccount - stock brokerage account in yaml file
type BrokerageAccount struct {
	Name    string  `yaml:"name"`
	Amount  float64 `yaml:"amount"`
	SaleFee float64 `yaml:"monthlySaleFee"`
}

// IRA in yaml file
type IRA struct {
	Name         string  `yaml:"name"`
	Amount       float64 `yaml:"amount"`
	PercentStock float64 `yaml:"percentStock"`
}

// SocialSecurity account in yaml file
type SocialSecurity struct {
	Early     float64 `yaml:"early"`
	Normal    float64 `yaml:"normal"`
	Late      float64 `yaml:"late"`
	Selection string  `yaml:"selection"`
}

// TaxBracket in yaml file
type TaxBracket struct {
	Minimum float64  `yaml:"minimum"`
	Maximum *float64 `yaml:"maximum"`
	Rate    float64  `yaml:"rate"`
}

// FederalBrackets in yaml file
type FederalBrackets struct {
	StandardDeduction float64      `yaml:"standardDeduction"`
	TaxBrackets       []TaxBracket `yaml:"brackets"`
}

// BracketCollection - yaml file struct containing all tax brackets
type BracketCollection struct {
	FederalBrackets      FederalBrackets `yaml:"federal"`
	StateBrackets        []TaxBracket    `yaml:"state"`
	CapitalGainsBrackets []TaxBracket    `yaml:"capitalGains"`
}

// YamlConfig - struct holding yaml file data, and data converted to model data for this application
type YamlConfig struct {
	YamlPerson         	   YamlPerson         `yaml:"person"`
	StockReturn            float64            `yaml:"stockReturn"`
	BondReturn         	   float64            `yaml:"bondReturn"`
	InflationRate      	   float64            `yaml:"inflationRate"`
	YamlCreditCards        []CreditCard       `yaml:"creditCards"`
	YamlInterestAccounts   []InterestAccount  `yaml:"interestAccounts"`
	YamlNoInterestAccounts []NoInterestAccount `yaml:"noInterestAccounts"`
	YamlBrokerageAccounts  []BrokerageAccount  `yaml:"brokerage"`
	YamlIRAs               []IRA               `yaml:"ira"`
	YamlSocialSecurity     SocialSecurity     `yaml:"socialSecurity"`
	YamlBracketCollection  BracketCollection  `yaml:"taxBrackets"`
	LogLevel 			   *string 			  `yaml:"logLevel"`
	MonthlyIncome 		   *float64 	      `yaml:"monthlyIncome"`
	person 				   *models.Person
	creditCards 		   models.CreditCards
	interestAccounts 	   models.AccountsWithInterest
	noInterestAccounts     models.AccountsNoInterest
	brokerageAccounts 	   models.AccountsStockBrokerage
	iras 				   models.IRAs
	socialSecurity 		   *models.AccountSocialSecurity
	bracketCollection 	   *models.TaxBracketCollection
}

// Load loads the config file into the yaml structure
func Load(file string) *YamlConfig {
	configData, err := os.ReadFile(file)
	if err != nil {
		logger.Get().Crit(fmt.Sprintf("Failed to read config file: %v", err))
	}

	var config YamlConfig
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		logger.Get().Crit(fmt.Sprintf("Failed to parse config file: %v", err))
	}
	logger.Get().Debug(fmt.Sprintf("%+v", config))
	return &config
}

// Person retrieves the person's data from the yaml file
func (y *YamlConfig) Person() *models.Person {
	if y.person == nil {
		y.person = models.NewPerson(y.YamlPerson.Years, y.YamlPerson.Months, y.YamlPerson.LifeExpectancy, y.YamlPerson.NeededPerMonth, y.InflationRate)
	}
	return y.person
}

// CreditCards retrieves credit card data from the yaml file
func (y *YamlConfig) CreditCards() models.CreditCards {
	if y.creditCards == nil {
		y.creditCards = models.CreditCards{}
		for _, yc := range y.YamlCreditCards {
			creditCard := models.NewCreditCardAccount(yc.Name, yc.Amount)
			y.creditCards = append(y.creditCards, creditCard)
		}
	}
	return y.creditCards
}

// InterestAccounts retrives passively-increasing cash account data from the yaml file
func (y *YamlConfig) InterestAccounts() models.AccountsWithInterest {
	if y.interestAccounts == nil {
		y.interestAccounts = models.AccountsWithInterest{}
		for _, ya := range y.YamlInterestAccounts {
			interestAccount := models.NewAccountWithInterest(ya.Name, ya.Amount, ya.Rate)
			y.interestAccounts = append(y.interestAccounts, interestAccount)
		}
	}
	return y.interestAccounts
}

// NoInterestAccounts retreives non-passively-increasing cash account data from the yaml file
func (y *YamlConfig) NoInterestAccounts() models.AccountsNoInterest {
	if y.noInterestAccounts == nil {
		y.noInterestAccounts = models.AccountsNoInterest{}
		for _, ya := range y.YamlNoInterestAccounts {
			noInterestAccount := models.NewAccountNoInterest(ya.Name, ya.Amount)
			y.noInterestAccounts = append(y.noInterestAccounts, noInterestAccount)
		}
	}
	return y.noInterestAccounts
}

// BrokerageAccounts retrieves brokerage account data from the yaml file
func (y *YamlConfig) BrokerageAccounts() models.AccountsStockBrokerage {
	if y.brokerageAccounts == nil {
		y.brokerageAccounts = models.AccountsStockBrokerage{}
		for _, ya := range y.YamlBrokerageAccounts {
			brokerageAccount := models.NewAccountStockBrokerage(ya.Name, ya.Amount, y.StockReturn, ya.SaleFee, y.Person())
			y.brokerageAccounts = append(y.brokerageAccounts, brokerageAccount)
		}
	}
	return y.brokerageAccounts
}

// IRAs retrieves IRA account data from the yaml file
func (y *YamlConfig) IRAs() models.IRAs {
	if y.iras == nil {
		y.iras = models.IRAs{}
		for _, ya := range y.YamlIRAs {
			ira := models.NewIRA(ya.Name, ya.Amount, ya.PercentStock, y.StockReturn, y.BondReturn, y.Person())
			y.iras = append(y.iras, ira)
		}
	}
	return y.iras
}

// SocialSecurity retrieves social security data from the yaml file
func (y *YamlConfig) SocialSecurity() *models.AccountSocialSecurity {
	if y.socialSecurity == nil {
		switch y.YamlSocialSecurity.Selection {
		case SSPayoutEarly:
			y.socialSecurity = models.NewAccountSocialSecurity(models.Early, y.YamlSocialSecurity.Early, y.InflationRate, y.Person())
		case SSPayoutNormal:
			y.socialSecurity = models.NewAccountSocialSecurity(models.Normal, y.YamlSocialSecurity.Normal, y.InflationRate, y.Person())
		case SSPayoutLate:
			y.socialSecurity = models.NewAccountSocialSecurity(models.Late, y.YamlSocialSecurity.Late, y.InflationRate, y.Person())
		default:
			panic(fmt.Sprintf("Invalid social security payout time: %s", y.YamlSocialSecurity.Selection))
		}
	}
	return y.socialSecurity
}

// TaxBracketCollection retrieves all tax bracket data from the yaml file
func (y *YamlConfig) TaxBracketCollection() *models.TaxBracketCollection {
	if y.bracketCollection == nil {
		ftb := []*models.TaxBracket{}
		for _, b := range y.YamlBracketCollection.FederalBrackets.TaxBrackets {
			if b.Maximum == nil {
				ftb = append(ftb, models.NewTaxBracketNoMax(b.Minimum, b.Rate))
				continue
			}
			ftb = append(ftb, models.NewTaxBracket(b.Minimum, *b.Maximum, b.Rate))
		}
		federalBrackets := models.NewFedTaxBrackets(y.YamlBracketCollection.FederalBrackets.StandardDeduction, ftb)
		stb := []*models.TaxBracket{}
		for _, b := range y.YamlBracketCollection.StateBrackets {
			if b.Maximum == nil {
				stb = append(stb, models.NewTaxBracketNoMax(b.Minimum, b.Rate))
				continue
			}
			stb = append(stb, models.NewTaxBracket(b.Minimum, *b.Maximum, b.Rate))
		}
		stateBrackets := models.NewStateTaxBrackets(stb)
		cgtb := []*models.TaxBracket{}
		for _, b := range y.YamlBracketCollection.CapitalGainsBrackets {
			if b.Maximum == nil {
				cgtb = append(cgtb, models.NewTaxBracketNoMax(b.Minimum, b.Rate))
				continue
			}
			cgtb = append(cgtb, models.NewTaxBracket(b.Minimum, *b.Maximum, b.Rate))
		}
		capGainsBrackets := models.NewCapTaxBrackets(cgtb)
		y.bracketCollection = models.NewTaxBracketCollection(federalBrackets, stateBrackets, capGainsBrackets)
	}
	return y.bracketCollection
}