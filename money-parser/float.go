package moneyparser

import (
	"math"

	"github.com/Rhymond/go-money"
)

// NewFromFloat creates and returns new instance of Money from a float64.
// Always rounding trailing decimals up.
func NewFromFloat(amount float64, currency string) *money.Money {
	currencyDecimals := math.Pow10(money.GetCurrency(currency).Fraction)
	ratio := math.Pow(10, 2)
	roundedAmount := math.Round(amount*ratio) / ratio
	return money.New(int64(roundedAmount*currencyDecimals), currency)
}
