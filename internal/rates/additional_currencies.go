package rates

import (
	"github.com/Rhymond/go-money"
)

func init() {
	// Some trading pairs require non-standard currencies that are not included in the go-money package by default.
	money.AddCurrency("BTC", "₿", "1 $", ".", ",", 8)
	money.AddCurrency("CNH", "¥", "1 $", ".", ",", 2)
	money.AddCurrency("XPD", "XPD", "1 $", ".", ",", 2)
	money.AddCurrency("XPT", "XPT", "1 $", ".", ",", 2)

	// Adding non standard currencies used in FixedCryptoProvider.
	money.AddCurrency("BEER", "BEER", "1 $", ".", ",", 18)
	money.AddCurrency("FLOKI", "FLOKI", "1 $", ".", ",", 18)
	money.AddCurrency("GATE", "GATE", "1 $", ".", ",", 18)
	money.AddCurrency("USDT", "USDT", "1 $", ".", ",", 6)
	money.AddCurrency("WBTC", "WBTC", "1 $", ".", ",", 8)

}
