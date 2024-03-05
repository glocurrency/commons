package moneyparser_test

import (
	"strconv"
	"testing"

	"github.com/Rhymond/go-money"
	moneyparser "github.com/glocurrency/commons/money-parser"
	"github.com/stretchr/testify/assert"
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
	}

	for i := range tests {
		test := tests[i]
		t.Run(strconv.FormatFloat(test.float, 'f', 0, 64), func(t *testing.T) {
			t.Parallel()

			m := moneyparser.NewFromFloat(test.float, test.currency)
			assert.Equal(t, test.wantCents, m.Amount())
			assert.Equal(t, test.currency, m.Currency().Code)
		})
	}
}
