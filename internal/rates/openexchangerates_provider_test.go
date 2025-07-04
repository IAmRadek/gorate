package rates

import (
	"net/http"
	"os"
	"testing"

	"github.com/Rhymond/go-money"
)

func TestOpenExchangeRatesProvider(t *testing.T) {
	prov := NewOpenExchangeRatesProvider(
		http.DefaultClient,
		os.Getenv("OPEN_EXCHANGE_RATES_PROVIDER_APP_ID"),
	)

	rates, err := prov.Rates(t.Context(),
		money.GetCurrency("USD"),
		money.GetCurrency("GBP"),
		money.GetCurrency("EUR"),
		money.GetCurrency("BTC"),
	)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expectedCombinations := [][2]string{
		{"USD", "BTC"},
		{"EUR", "BTC"},
		{"USD", "BTC"},
		{"GBP", "BTC"},
		{"USD", "EUR"},
		{"BTC", "EUR"},
		{"USD", "EUR"},
		{"GBP", "EUR"},
		{"USD", "GBP"},
		{"BTC", "GBP"},
		{"USD", "GBP"},
		{"EUR", "GBP"},
		{"EUR", "USD"},
		{"GBP", "USD"},
		{"GBP", "USD"},
		{"EUR", "USD"},
		{"BTC", "USD"},
		{"BTC", "USD"},
	}

	if len(rates) != len(expectedCombinations) {
		t.Fatalf("Expected %d rates go %d", len(expectedCombinations), len(rates))
	}

	for i, ec := range expectedCombinations {
		if ec[0] != rates[i].From.Code {
			t.Fatalf("Expected %d From currency to match %q got %q", i, ec[0], rates[i].From.Code)
		}

		if ec[1] != rates[i].To.Code {
			t.Fatalf("Expected %d To currency to match %q got %q", i, ec[0], rates[i].From.Code)
		}
	}
}
