package main

import (
	"github.com/lrm25/moneyLeft/application"
	"github.com/lrm25/moneyLeft/config"
	"github.com/lrm25/moneyLeft/logger"
)

func main() {

	logger.Init(logger.LEVEL_TRACE)
	c := config.Load("accounts.yaml")
	application.Run(c)
}