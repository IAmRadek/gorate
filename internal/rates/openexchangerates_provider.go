package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/Rhymond/go-money"
	"github.com/govalues/decimal"
)

type OpenExchangeRatesProvider struct {
	client *http.Client
	appID  string

	// TODO: implement caching.
}

func NewOpenExchangeRatesProvider(cli *http.Client, appID string) *OpenExchangeRatesProvider {
	return &OpenExchangeRatesProvider{
		client: cli,
		appID:  appID,
	}
}

func (o *OpenExchangeRatesProvider) SupportedCurrencies(ctx context.Context) ([]*money.Currency, error) {
	currencies, err := o.getCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting currencies: %w", err)
	}

	return currencies, nil
}

func (o *OpenExchangeRatesProvider) Rates(ctx context.Context, c1, c2 *money.Currency, c ...*money.Currency) (ExchangeRates, error) {
	currencies := []*money.Currency{c1, c2}
	currencies = append(currencies, c...)

	rates, err := o.getRates(ctx, currencies)
	if err != nil {
		return nil, fmt.Errorf("getting rates: %w", err)
	}

	out := make([]ExchangeRate, 0, len(rates))

	for _, rate := range rates {
		expectedFrom := slices.ContainsFunc(currencies, func(currency *money.Currency) bool {
			return currency.Code == rate.From.Code
		})

		expectedTo := slices.ContainsFunc(currencies, func(currency *money.Currency) bool {
			return currency.Code == rate.To.Code
		})

		if !expectedFrom || !expectedTo {
			continue
		}

		out = append(out, rate)
	}

	return out, nil

}

// getRates retrieves exchange rates for the provided currencies using the Open Exchange Rates API.
// It ensures there are at least two distinct currencies and computes cross-rates for all currency pairs.
// Returns a list of ExchangeRate containing rate information or an error if the retrieval or processing fails.
func (o *OpenExchangeRatesProvider) getRates(ctx context.Context, currencies []*money.Currency) ([]ExchangeRate, error) {
	uniq := map[string]struct{}{}
	for _, s := range currencies {
		uniq[s.Code] = struct{}{}
	}
	if len(uniq) < 2 {
		return nil, fmt.Errorf("at least 2 distinct currencies required")
	}

	reqSymbols := make([]string, 0, len(uniq))
	for c := range uniq {
		if c == money.USD {
			continue
		}
		reqSymbols = append(reqSymbols, c)
	}
	sort.Strings(reqSymbols)

	params := url.Values{}
	params.Add("app_id", o.appID)
	params.Add("symbols", strings.Join(reqSymbols, ","))
	url := "https://openexchangerates.org/api/latest.json?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	var raw struct {
		Rates map[string]decimal.Decimal `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding data: %w", err)
	}

	if len(raw.Rates) == 0 {
		return nil, fmt.Errorf("no rates found")
	}

	raw.Rates[money.USD] = decimal.One

	for c := range uniq {
		if _, ok := raw.Rates[c]; !ok {
			return nil, fmt.Errorf("openexchangerates missing rate for %q", c)
		}
	}

	currList := make([]string, 0, len(uniq))
	for c := range uniq {
		currList = append(currList, c)
	}
	sort.Strings(currList)

	usd := money.GetCurrency(money.USD)
	out := make([]ExchangeRate, 0, len(currList)*len(currList))

	for _, from := range currList {
		curFrom := money.GetCurrency(from)
		if curFrom == nil {
			return nil, fmt.Errorf("unknown currency: %q", from)
		}
		rateFrom := raw.Rates[from]

		for _, to := range currList {
			if from == to {
				continue
			}
			curTo := money.GetCurrency(to)
			if curTo == nil {
				return nil, fmt.Errorf("unknown currency: %q", to)
			}
			rateTo := raw.Rates[to]

			cross, err := rateTo.Quo(rateFrom)
			if err != nil {
				return nil, fmt.Errorf("making cross rate: %w", err)
			}

			out = append(out, ExchangeRate{
				From: curFrom,
				To:   curTo,
				Rate: cross,
			})
		}

		if from != money.USD {
			r, _ := decimal.One.Quo(rateFrom)

			out = append(out,
				ExchangeRate{From: curFrom, To: usd, Rate: r},
				ExchangeRate{From: usd, To: curFrom, Rate: rateFrom},
			)
		}
	}

	slices.SortFunc(out, func(a, b ExchangeRate) int {
		return strings.Compare(a.To.Code, b.To.Code)
	})

	return out, nil
}

func (o *OpenExchangeRatesProvider) getCurrencies(ctx context.Context) ([]*money.Currency, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openexchangerates.org/api/currencies.json", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	currencies := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		return nil, fmt.Errorf("decoding data: %w", err)
	}

	out := make([]*money.Currency, 0, len(currencies))
	for code := range currencies {
		currency := money.GetCurrency(code)
		if currency == nil {
			slog.ErrorContext(ctx, "currency not found", "currency", code)
			continue
		}

		out = append(out, currency)
	}

	slices.SortFunc(out, func(a, b *money.Currency) int {
		return strings.Compare(a.Code, b.Code)
	})

	return out, nil
}
