package main

import (
	"fmt"
	"net/http"

	"github.com/IAmRadek/gorate/internal/exchanges"
	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
	"github.com/govalues/decimal"
)

func HandleExchange(exchange *exchanges.Exchange) gin.HandlerFunc {
	type request struct {
		From   string          `form:"from"`
		To     string          `form:"to"`
		Amount decimal.Decimal `form:"amount"`
	}

	type response struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	return func(c *gin.Context) {
		var req request

		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("cannot parse request: %v", err),
			})
			return
		}

		from := money.GetCurrency(req.From)
		if from == nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		to := money.GetCurrency(req.To)
		if to == nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		if req.Amount.Sign() <= 0 {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		m, err := exchange.Exchange(c.Copy(), from, to, req.Amount)
		if err != nil {
			resp := map[string]any{
				"error": fmt.Sprintf("exchange failed: %v", err),
			}

			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		c.JSON(http.StatusOK, response{
			From:   from.Code,
			To:     to.Code,
			Amount: m.AsMajorUnits(),
		})

	}
}

func filter(su []string, code string) []string {
	filtered := make([]string, 0, len(su))
	for _, s := range su {
		if s != code {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
