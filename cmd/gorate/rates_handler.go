package main

import (
	"net/http"
	"strings"

	"github.com/IAmRadek/gorate/internal/rates"
	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

func HandleRates(rates rates.Provider) gin.HandlerFunc {
	type request struct {
		Currencies string `form:"currencies"`
	}

	type response struct {
		From string  `json:"from,omitempty"`
		To   string  `json:"to,omitempty"`
		Rate float64 `json:"rate,omitempty"`
	}

	return func(c *gin.Context) {
		var req request

		if err := c.ShouldBindQuery(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		rawCurrencies := strings.Split(req.Currencies, ",")

		if len(rawCurrencies) < 2 {
			c.Status(http.StatusBadRequest)
			return
		}

		currencies := make([]*money.Currency, 0, len(req.Currencies))
		for _, cur := range rawCurrencies {
			currency := money.GetCurrency(cur)
			if currency == nil {
				c.Status(http.StatusBadRequest)
				return
			}

			currencies = append(currencies, currency)
		}

		exchangeRates, err := rates.Rates(c.Copy(), currencies[0], currencies[1], currencies[1:]...)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out := make([]response, 0, len(exchangeRates))

		for _, rate := range exchangeRates {
			r, _ := rate.Rate.Float64()

			out = append(out, response{
				From: rate.From.Code,
				To:   rate.To.Code,
				Rate: r,
			})
		}

		c.JSON(http.StatusOK, out)
	}
}
