package moneyparser_test

import (
	"strconv"
	"testing"

	"github.com/Rhymond/go-money"
	moneyparser "github.com/glocurrency/commons/money-parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromFloat(t *testing.T) {
	tests := []struct {
		float     float64
		currency  string
		wantCents int64
	}{
		{12.34, money.EUR, 1234},
		{-0.125, money.EUR, -13},
		{-0.126, money.EUR, -13},
		{0.129, money.EUR, 13},
		{9823.21, money.EUR, 982321},
		{1.13, money.EUR, 113},
		{39.99, money.EUR, 3999},
	}

	for i := range tests {
		test := tests[i]
		t.Run(strconv.FormatFloat(test.float, 'f', -1, 64), func(t *testing.T) {
			t.Parallel()

			m := moneyparser.NewFromFloat(test.float, test.currency)
			assert.Equal(t, test.wantCents, m.Amount())
			assert.Equal(t, test.currency, m.Currency().Code)
		})
	}
}

func TestNewFromFloat_MathToZero(t *testing.T) {
	// The specific sequence that caused the "-0.01" bug due to float drift
	amounts := []float64{9823.21, -9400.00, -400.00, -23.21}

	balance := money.New(0, money.EUR)
	var err error

	for _, amount := range amounts {
		parsedAmount := moneyparser.NewFromFloat(amount, money.EUR)
		balance, err = balance.Add(parsedAmount)
		require.NoError(t, err)
	}

	// If precision was lost during parsing, this would be non-zero (e.g. -1)
	assert.Equal(t, int64(0), balance.Amount())
}
