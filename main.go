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
	c := config.Load("accounts.yaml")
	if c.LogLevel != nil {
		logger.Get().SetLevel(*c.LogLevel)
	}

	// Define a command line flag named "number"
	costToConvert := flag.Float64("costToConvert", 0.0, "Cost to convert to time")
	flag.Parse()
	if 0.01 <= *costToConvert {
		printCostToConvert(*costToConvert, c.YamlPerson.NeededPerMonth)
		return
	}
	application.Run(c)
}