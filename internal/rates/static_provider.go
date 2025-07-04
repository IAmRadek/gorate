package rates

import (
	"context"

	"github.com/Rhymond/go-money"
	"github.com/govalues/decimal"
)

type StaticTestRatesProvider struct{}

func NewStaticRatesProvider() StaticTestRatesProvider {
	return StaticTestRatesProvider{}
}

func (s StaticTestRatesProvider) SupportedCurrencies(ctx context.Context) ([]*money.Currency, error) {
	return []*money.Currency{
		money.GetCurrency("USD"),
		money.GetCurrency("GBP"),
		money.GetCurrency("EUR"),
		money.GetCurrency("BTC"),
	}, nil
}

func (s StaticTestRatesProvider) Rates(ctx context.Context, c1, c2 *money.Currency, c ...*money.Currency) (ExchangeRates, error) {
	return []ExchangeRate{
		{money.GetCurrency("USD"), money.GetCurrency("BTC"), decimal.MustParse("0.000009104837")},
		{money.GetCurrency("EUR"), money.GetCurrency("BTC"), decimal.MustParse("0.0000106959819745101")},
		{money.GetCurrency("USD"), money.GetCurrency("BTC"), decimal.MustParse("0.000009104837")},
		{money.GetCurrency("GBP"), money.GetCurrency("BTC"), decimal.MustParse("0.0000124249434010156")},
		{money.GetCurrency("USD"), money.GetCurrency("EUR"), decimal.MustParse("0.851239")},
		{money.GetCurrency("BTC"), money.GetCurrency("EUR"), decimal.MustParse("93493.05209966965911")},
		{money.GetCurrency("USD"), money.GetCurrency("EUR"), decimal.MustParse("0.851239")},
		{money.GetCurrency("GBP"), money.GetCurrency("EUR"), decimal.MustParse("1.161645880726595859")},
		{money.GetCurrency("USD"), money.GetCurrency("GBP"), decimal.MustParse("0.732787")},
		{money.GetCurrency("BTC"), money.GetCurrency("GBP"), decimal.MustParse("80483.26400571476458")},
		{money.GetCurrency("USD"), money.GetCurrency("GBP"), decimal.MustParse("0.732787")},
		{money.GetCurrency("EUR"), money.GetCurrency("GBP"), decimal.MustParse("0.860847541054862383")},
		{money.GetCurrency("EUR"), money.GetCurrency("USD"), decimal.MustParse("1.174758205392375114")},
		{money.GetCurrency("GBP"), money.GetCurrency("USD"), decimal.MustParse("1.364653030143820783")},
		{money.GetCurrency("GBP"), money.GetCurrency("USD"), decimal.MustParse("1.364653030143820783")},
		{money.GetCurrency("EUR"), money.GetCurrency("USD"), decimal.MustParse("1.174758205392375114")},
		{money.GetCurrency("BTC"), money.GetCurrency("USD"), decimal.MustParse("109831.7301012637568")},
		{money.GetCurrency("BTC"), money.GetCurrency("USD"), decimal.MustParse("109831.7301012637568")},
	}, nil
}
