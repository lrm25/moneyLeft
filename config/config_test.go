package config

import (
	"os"
	"testing"

	"github.com/lrm25/moneyLeft/models"
	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	testBytes := []byte(
`person:
  years: 50
  months: 5
  lifeExpectancy: 70
  neededPerMonth: 3000.00
stockReturn: 10.0
bondReturn: 5.5
inflationRate: 3.1
creditCards:
  - name: ABCard
    amount: 2525.25
  - name: DEF Bank
    amount: 1500.00
interestAccounts:
  - name: Test Bank
    amount: 10000.0
    rate: 2.25
noInterestAccounts:
  - name: Boring Bank
    amount: 5000.00
brokerage:
  - name: Some Stock Company
    amount: 25000.00
    monthlySaleFee: 9.99
ira:
  - name: IRA Company
    amount: 50000.00
    percentStock: 50
socialSecurity:
  early: 1000
  normal: 2000
  late: 3000
  selection: normal
taxBrackets:
  federal:
    standardDeduction: 13850
    brackets:
    - minimum: 0
      maximum: 11000
      rate: 10
    - minimum: 11000
      maximum: 44725
      rate: 12
  state:
    - minimum: 0
      rate: 4.4
  capitalGains:
    - minimum: 0
      maximum: 44625
      rate: 0
    - minimum: 44625
      maximum: 492300
      rate: 15`)
	configFileName := "testConfig.yaml"
	testConfig, err := os.Create(configFileName)
	require.NoError(t, err)
	defer func() {
		testConfig.Close()
		os.Remove(configFileName)
	}()

	_, err = testConfig.Write(testBytes)
	require.NoError(t, err)

	c := Load(configFileName)
	person := c.Person()
	require.Equal(t, 50, person.AgeYears())
	require.Equal(t, 5, person.AgeMonths())
	require.Equal(t, 70, person.LifeExpectancy())
	require.InDelta(t, 3000.00, person.NeededPerMonth(), 0.01)
	require.InDelta(t, 10.0, c.StockReturn, 0.01)
	require.InDelta(t, 5.5, c.BondReturn, 0.01)
	require.InDelta(t, 3.1, c.InflationRate, 0.01)

	creditCards := c.CreditCards()
	require.Equal(t, 2, len(creditCards))
	require.Equal(t, creditCards[0].Name(), "ABCard")
	require.InDelta(t, 2525.25, creditCards[0].Amount(), 0.01)
	require.Equal(t, creditCards[1].Name(), "DEF Bank")
	require.InDelta(t, 1500.00, creditCards[1].Amount(), 0.01)

	interestAccounts := c.InterestAccounts()
	require.Equal(t, 1, len(interestAccounts))
	require.InDelta(t, 10000.00, interestAccounts[0].Amount(), 0.01)
	require.InDelta(t, 2.25, interestAccounts[0].Rate(), 0.01)
	
	noInterestAccounts := c.NoInterestAccounts()
	require.Equal(t, 1, len(noInterestAccounts))
	require.InDelta(t, 5000.00, noInterestAccounts[0].Amount(), 0.01)

	brokerageAccounts := c.BrokerageAccounts()
	require.Equal(t, 1, len(brokerageAccounts))
	require.Equal(t, "Some Stock Company", brokerageAccounts[0].Name())
	require.InDelta(t, 25000.00, brokerageAccounts[0].Amount(), 0.01)
	require.InDelta(t, 9.99, brokerageAccounts[0].MonthlySaleFee(), 0.01)
	require.Equal(t, c.StockReturn, brokerageAccounts[0].Rate())
	require.Equal(t, 50, brokerageAccounts[0].Person().AgeYears())

	iraAccounts := c.IRAs()
	require.Equal(t, 1, len(iraAccounts))
	require.Equal(t, "IRA Company", iraAccounts[0].Name())
	require.InDelta(t, 50000.00, iraAccounts[0].Amount(), 0.01)
	require.InDelta(t, 50.0, iraAccounts[0].PercentStock(), 0.01)
	require.InDelta(t, c.StockReturn, iraAccounts[0].StockInterestRate(), 0.01)
	require.InDelta(t, c.BondReturn, iraAccounts[0].BondInterestRate(), 0.01)
	require.Equal(t, 50, iraAccounts[0].Person().AgeYears())
	
	socialSecurity := c.SocialSecurity()
	// not 67 yet
	require.InDelta(t, 0.0, socialSecurity.Amount(), 0.01)
	require.Equal(t, models.NORMAL, socialSecurity.Selection())
	require.Equal(t, 50, socialSecurity.Person().AgeYears())

	taxBrackets := c.TaxBracketCollection()
	fedBrackets := taxBrackets.Federal()
	require.Equal(t, 13850.00, fedBrackets.StandardDeduction())
	require.Equal(t, 2, len(fedBrackets.Brackets()))
	require.InDelta(t, 0.00, fedBrackets.Brackets()[0].Minimum(), 0.01)
	require.InDelta(t, 11000.00, fedBrackets.Brackets()[0].Maximum(), 0.01)
	require.InDelta(t, 11000.00, fedBrackets.Brackets()[1].Minimum(), 0.01)
	require.InDelta(t, 44725.00, fedBrackets.Brackets()[1].Maximum(), 0.01)
	require.Equal(t, true, fedBrackets.Brackets()[0].HasMaximum())
	stateBrackets := taxBrackets.State()
	require.Equal(t, 1, len(stateBrackets.Brackets()))
	require.InDelta(t, 0.0, stateBrackets.Brackets()[0].Minimum(), 0.01)
	require.Equal(t, false, stateBrackets.Brackets()[0].HasMaximum())
	capTaxBrackets := taxBrackets.CapitalGains()
	require.Equal(t, 2, len(capTaxBrackets.Brackets()))
	require.InDelta(t, 0.0, capTaxBrackets.Brackets()[0].Minimum(), 0.01)
	require.InDelta(t, 44625.00, capTaxBrackets.Brackets()[0].Maximum(), 0.01)
	require.InDelta(t, 492300.00, capTaxBrackets.Brackets()[1].Maximum(), 0.01)
}