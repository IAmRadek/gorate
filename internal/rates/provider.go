package rates

import (
	"context"

	"github.com/Rhymond/go-money"
)

// Provider encapsulates different rates providers that may support different Currencies and be refreshed at different rates.
type Provider interface {
	// SupportedCurrencies returns a list of currencies that we can expect from Provider.
	SupportedCurrencies(ctx context.Context) ([]*money.Currency, error)

	// Rates returns current rates for a given set of currencies at least two is required.
	Rates(ctx context.Context, c1, c2 *money.Currency, c ...*money.Currency) (ExchangeRates, error)
}
