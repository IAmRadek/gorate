package exchanges

import (
	"context"
	"fmt"

	"github.com/IAmRadek/gorate/internal/rates"
	"github.com/Rhymond/go-money"
	"github.com/govalues/decimal"
)

type Exchange struct {
	provider rates.Provider
}

func NewExchange(prov rates.Provider) *Exchange {
	return &Exchange{
		provider: prov,
	}
}

func (ex *Exchange) SupportedCurrencies(ctx context.Context) ([]string, error) {
	sc, err := ex.provider.SupportedCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting supported currencies: %w", err)
	}

	out := make([]string, 0, len(sc))
	for _, s := range sc {
		out = append(out, s.Code)
	}

	return out, err
}

func (ex *Exchange) Exchange(ctx context.Context, from, to *money.Currency, amount decimal.Decimal) (*money.Money, error) {
	rates, err := ex.provider.Rates(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("getting rates for %q and %q: %w", from.Code, to.Code, err)
	}

	rate, found := rates.For(from, to)
	if !found {
		return nil, fmt.Errorf("rate for %q and %q is not possible", from.Code, to.Code)
	}

	// NOTE: in here we could also insert an external component for adding additional fees etc.
	newAmount, err := amount.Mul(rate.Rate)
	if err != nil {
		return nil, fmt.Errorf("calculating new amount: %w", err)
	}

	fl, ok := newAmount.Float64()
	if !ok {
		return nil, fmt.Errorf("invalid floating point: %q", newAmount.String())
	}

	return money.NewFromFloat(fl, to.Code), nil
}
