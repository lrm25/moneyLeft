package main

import (
	"flag"
	"fmt"

	"github.com/lrm25/moneyLeft/application"
	"github.com/lrm25/moneyLeft/config"
	"github.com/lrm25/moneyLeft/logger"
)

func printCostToConvert(costToConvert, neededPerMonth float64) {
	months := costToConvert / neededPerMonth
	logger.Get().Info(fmt.Sprintf("Cost in months: %.02f, approx days: %.02f", months, months*30))
}

func main() {

	logger.Init(logger.LevelInfo)
	costToConvert := flag.Float64("costToConvert", 0.0, "Cost to convert to time")
	calculateMinNeeded := flag.Bool("calculateMinNeeded", false, "If true, calculated minimum needed to not go broke lifetime")
	flag.Parse()
	monthlyIncome := 0.00
	broke := true
	for broke {
		c := config.Load("accounts.yaml")
		if *calculateMinNeeded {
			c.MonthlyIncome = &monthlyIncome
		}
		if c.LogLevel != nil {
			logger.Get().SetLevel(*c.LogLevel)
		}

		if 0.01 <= *costToConvert {
			printCostToConvert(*costToConvert, c.YamlPerson.NeededPerMonth)
			return
		}
		broke = application.Run(c, *calculateMinNeeded)
		if !*calculateMinNeeded {
			break
		}
		if broke {
			monthlyIncome += 1.00
		}
	}
	if *calculateMinNeeded {
		fmt.Printf("income required: %.2f\n", monthlyIncome)
	}
}
