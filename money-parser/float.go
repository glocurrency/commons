package moneyparser

import (
	"math"

	"github.com/Rhymond/go-money"
)

// NewFromFloat creates and returns new instance of Money from a float64.
// Always rounding trailing decimals up.
func NewFromFloat(amount float64, currency string) *money.Money {
	currencyDecimals := math.Pow10(money.GetCurrency(currency).Fraction)
	minorUnits := int64(math.Round(amount * currencyDecimals))
	return money.New(minorUnits, currency)
}
