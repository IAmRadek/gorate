package e2e_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type response struct {
	From string  `json:"from,omitempty"`
	To   string  `json:"to,omitempty"`
	Rate float64 `json:"rate,omitempty"`
}

func TestRatesE2E(t *testing.T) {
	baseURL := "http://localhost:8080"

	// Check if the server is already running
	resp, err := http.Get(baseURL + "/rates")
	if err != nil {
		t.Logf("Pinging server error: %v", err)
		t.Skip("Server is not running. Start the server before running this test.")
	}
	resp.Body.Close()

	tests := []struct {
		name                string
		requestedCurrencies []string
		expectedCode        int
		expectedResponse    any
	}{
		{
			name:                "single_currency",
			requestedCurrencies: []string{"USD"},
			expectedCode:        http.StatusBadRequest,
			expectedResponse:    nil,
		},
		{
			name:                "unknown_currency",
			requestedCurrencies: []string{"ABC", "CDE"},
			expectedCode:        http.StatusBadRequest,
			expectedResponse:    []byte{},
		},
		{
			name:                "USD,GBP,EUR",
			requestedCurrencies: []string{"USD", "GBP", "EUR"},
			expectedCode:        http.StatusOK,
			expectedResponse: []response{
				{"USD", "EUR", 0.848818},
				{"GBP", "EUR", 1.1608418386535178},
				{"USD", "EUR", 0.848818},
				{"EUR", "GBP", 0.8614437959609716},
				{"USD", "GBP", 0.731209},
				{"USD", "GBP", 0.731209},
				{"EUR", "USD", 1.1781088525455399},
				{"EUR", "USD", 1.1781088525455399},
				{"GBP", "USD", 1.3675980465229503},
				{"GBP", "USD", 1.3675980465229503},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			vals := url.Values{}
			vals.Add("currencies", strings.Join(tc.requestedCurrencies, ","))

			resp, err := http.Get(baseURL + "/rates?" + vals.Encode())
			if err != nil {
				t.Fatalf("calling rates: %v", err)
			}

			if resp.StatusCode != tc.expectedCode {
				t.Fatalf("expected: %d status got: %d", tc.expectedCode, resp.StatusCode)
			}

			buf, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			var respBody any
			switch tc.expectedResponse.(type) {
			case nil:
				if len(buf) != 0 {
					t.Fatalf("Expected empty response got: %s", buf)
				}
			case []response:
				var rateResp []response
				if err := json.Unmarshal(buf, &rateResp); err != nil {
					t.Fatalf("decoding rates response: %v", err)
				}
				respBody = rateResp

				diff := cmp.Diff(respBody, tc.expectedResponse, cmpopts.EquateApprox(relTol, absTol))
				if diff != "" {
					t.Errorf("response mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

const (
	relTol = 1e-9  // ~9 decimal places
	absTol = 1e-12 // good down near zero
)
