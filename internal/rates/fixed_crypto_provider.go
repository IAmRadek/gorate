package rates

import (
	"context"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/govalues/decimal"
)

type FixedCryptoRatesProvider struct{}

func NewFixedCryptoRatesProvider() FixedCryptoRatesProvider {
	return FixedCryptoRatesProvider{}
}

func (s FixedCryptoRatesProvider) SupportedCurrencies(ctx context.Context) ([]*money.Currency, error) {
	return []*money.Currency{
		money.GetCurrency("USD"),
		money.GetCurrency("BEER"),
		money.GetCurrency("FLOKI"),
		money.GetCurrency("GATE"),
		money.GetCurrency("USDT"),
		money.GetCurrency("WBTC"),
	}, nil
}

func (s FixedCryptoRatesProvider) Rates(ctx context.Context, c1, c2 *money.Currency, c ...*money.Currency) (ExchangeRates, error) {
	ratesToUSD := ExchangeRates{
		{
			From: money.GetCurrency("USD"),
			To:   money.GetCurrency("USD"),
			Rate: decimal.One,
		},
		{
			From: money.GetCurrency("BEER"),
			To:   money.GetCurrency("USD"),
			Rate: decimal.MustParse("0.00002461"),
		},
		{
			From: money.GetCurrency("FLOKI"),
			To:   money.GetCurrency("USD"),
			Rate: decimal.MustParse("0.0001428"),
		},
		{
			From: money.GetCurrency("GATE"),
			To:   money.GetCurrency("USD"),
			Rate: decimal.MustParse("6.87"),
		},
		{
			From: money.GetCurrency("USDT"),
			To:   money.GetCurrency("USD"),
			Rate: decimal.MustParse("0.999"),
		},
		{
			From: money.GetCurrency("WBTC"),
			To:   money.GetCurrency("USD"),
			Rate: decimal.MustParse("57037.22"),
		},
	}

	usd := money.GetCurrency(money.USD)
	currList := []string{
		"USD",
		"BEER",
		"FLOKI",
		"GATE",
		"USDT",
		"WBTC",
	}

	out := make([]ExchangeRate, 0, len(currList)*len(currList))
	for _, from := range currList {
		curFrom := money.GetCurrency(from)
		if curFrom == nil {
			return nil, fmt.Errorf("unknown currency: %q", from)
		}

		rateFrom, _ := ratesToUSD.For(curFrom, usd)

		for _, to := range currList {
			if from == to {
				continue
			}

			curTo := money.GetCurrency(to)
			if curTo == nil {
				return nil, fmt.Errorf("unknown currency: %q", to)
			}

			rateTo, _ := ratesToUSD.For(curTo, usd)

			cross, err := rateTo.Rate.Quo(rateFrom.Rate)
			if err != nil {
				return nil, fmt.Errorf("making cross rate for %q and %q: %w", from, to, err)
			}

			out = append(out, ExchangeRate{
				From: curFrom,
				To:   curTo,
				Rate: cross,
			})
		}

		if from != money.USD {
			r, _ := decimal.One.Quo(rateFrom.Rate)

			out = append(out,
				ExchangeRate{From: curFrom, To: usd, Rate: r},
				ExchangeRate{From: usd, To: curFrom, Rate: rateFrom.Rate},
			)
		}
	}

	return out, nil
}
