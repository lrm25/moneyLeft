package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIRA_Deduct_correctPenalties(t *testing.T) {
	person := NewPerson(58, 0, 70, 1000.00, 0.00)
	ira := NewIRA("name", 50000.00, 0.0, 0.0, 0.0, person)
	expectedAmount := 50000.00
	for person.AgeYears() <= 60 {
		if person.AgeYears() < 59 || (person.AgeYears() == 59 && person.AgeMonths() < 6) {
			expectedAmount -= 1100.00
		} else {
			expectedAmount -= 1000.00
		}
		ira.Deduct(1000.00)
		require.InDelta(t, expectedAmount, ira.Amount(), 0.01, "expected amount mismatch, expected '%.2f', was '%.2f'", expectedAmount, ira.Amount())
		person.IncreaseAge(0, 1)
	}
}
