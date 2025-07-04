package exchanges

import (
	"testing"

	"github.com/IAmRadek/gorate/internal/rates"
	"github.com/Rhymond/go-money"
	"github.com/govalues/decimal"
)

func TestExchange(t *testing.T) {
	exch := NewExchange(rates.NewStaticRatesProvider())

	usd := money.GetCurrency("USD")
	btc := money.GetCurrency("BTC")

	ex, err := exch.Exchange(t.Context(), usd, btc, decimal.MustParse("1"))
	if err != nil {
		return
	}

	if ex.IsZero() {
		t.Fatalf("Expected exchanged to be non zero")
	}
}
