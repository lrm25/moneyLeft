package application

import (
	"bytes"
	"testing"

	"github.com/lrm25/moneyLeft/logger"
	"github.com/lrm25/moneyLeft/models"
	"github.com/stretchr/testify/require"
)

func Test_runAgeLoop(t *testing.T) {

	var buffer bytes.Buffer
	logger.InitWithWriter(logger.LevelInfo, &buffer)

	person := models.NewPerson(50, 0, 60, 1001.00, 0.00)
	person.WithAccounts(nil, models.PositiveAccounts{models.NewAccountNoInterest("test", 1000.00)}, nil)
	runAgeLoop(person, 2000, 1)
	require.Contains(t, string(buffer.Bytes()), "Broke on year 2000, month 2 (age 50 years, 1 months)")
}