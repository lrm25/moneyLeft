package main

import (
	"bytes"
	"testing"

	"github.com/lrm25/moneyLeft/logger"
	"github.com/stretchr/testify/require"
)

func Test_printCostToConvert(t *testing.T) {

	var buffer bytes.Buffer

	logger.InitWithWriter(logger.LevelInfo, &buffer)
	printCostToConvert(100.00, 50.00)
	require.Contains(t, string(buffer.Bytes()), "Cost in months: 2.00, approx days: 60.00")
}
