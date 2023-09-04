package main

import (
	"github.com/lrm25/moneyLeft/application"
	"github.com/lrm25/moneyLeft/config"
	"github.com/lrm25/moneyLeft/logger"
)

func main() {

	logger.Init(logger.LEVEL_INFO)
	c := config.Load("accounts.yaml")
	if c.LogLevel != nil {
		logger.Get().SetLevel(*c.LogLevel)
	}
	application.Run(c)
}