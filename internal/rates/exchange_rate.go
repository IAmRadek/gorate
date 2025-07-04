package rates

import (
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/govalues/decimal"
)

// ExchangeRate represents an exchange rate between two currencies, including their conversion rate.
type ExchangeRate struct {
	// From represents the source currency in the exchange rate.
	From *money.Currency

	// To represents the target currency in the exchange rate.
	To *money.Currency

	// Rate represents the conversion rate between two currencies in an exchange rate.
	Rate decimal.Decimal
}

func (r ExchangeRate) String() string {
	return fmt.Sprintf("%q => %q (%v)", r.From.Code, r.To.Code, r.Rate.String())
}

type ExchangeRates []ExchangeRate

func (r ExchangeRates) For(from, to *money.Currency) (ExchangeRate, bool) {
	for _, rate := range r {
		if rate.From.Code == from.Code && rate.To.Code == to.Code {
			return rate, true
		}
	}

	return ExchangeRate{}, false
}
