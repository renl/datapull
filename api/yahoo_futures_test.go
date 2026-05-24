package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestYahooFinanceClassifier(t *testing.T) {
	t.Run("full ohlc maps latest usable bar", func(t *testing.T) {
		payload := loadYahooFinanceEnvelope(t, `{
			"chart": {
				"result": [{
					"meta": {"symbol": "GC=F"},
					"timestamp": [1710000000, 1710086400, 1710172800],
					"indicators": {
						"quote": [{
							"open": [2300.1, 2305.1, 2310.1],
							"high": [2310.0, 2315.0, 2320.0],
							"low": [2295.0, 2300.0, 2305.0],
							"close": [2308.0, null, 2318.0],
							"volume": [100, 110, 120]
						}]
					}
				}],
				"error": null
			}
		}`)

		classification, result, err := classifyYahooFinanceResponse(payload, YahooFinanceGoldFuturesSymbol)
		if err != nil {
			t.Fatalf("classifyYahooFinanceResponse() error = %v", err)
		}
		if classification != YahooFinanceFullOHLCResponse {
			t.Fatalf("classification = %s, want %s", classification, YahooFinanceFullOHLCResponse)
		}

		data, err := mapYahooFinanceCommodityCardData(result, classification)
		if err != nil {
			t.Fatalf("mapYahooFinanceCommodityCardData() error = %v", err)
		}

		if data.Close != 2318.0 || data.Price != 2318.0 {
			t.Fatalf("mapped close/price = %v/%v, want 2318", data.Close, data.Price)
		}
		if data.Date != "1710172800" {
			t.Fatalf("mapped date = %q, want 1710172800", data.Date)
		}
		if data.Open == nil || *data.Open != 2310.1 {
			t.Fatalf("mapped open = %v, want 2310.1", data.Open)
		}
		if data.High == nil || *data.High != 2320.0 {
			t.Fatalf("mapped high = %v, want 2320.0", data.High)
		}
		if data.Low == nil || *data.Low != 2305.0 {
			t.Fatalf("mapped low = %v, want 2305.0", data.Low)
		}
		if data.Volume == nil || *data.Volume != 120 {
			t.Fatalf("mapped volume = %v, want 120", data.Volume)
		}
	})

	t.Run("partial meta only maps fallback fields", func(t *testing.T) {
		payload := loadYahooFinanceEnvelope(t, `{
			"chart": {
				"result": [{
					"meta": {
						"symbol": "CL=F",
						"regularMarketPrice": 81.23,
						"regularMarketDayHigh": 82.10,
						"regularMarketDayLow": 80.90,
						"regularMarketVolume": 456789,
						"regularMarketTime": 1710259200,
						"chartPreviousClose": 80.75
					},
					"indicators": {
						"quote": [{}]
					}
				}],
				"error": null
			}
		}`)

		classification, result, err := classifyYahooFinanceResponse(payload, YahooFinanceOilFuturesSymbol)
		if err != nil {
			t.Fatalf("classifyYahooFinanceResponse() error = %v", err)
		}
		if classification != YahooFinancePartialMetaOnlyResponse {
			t.Fatalf("classification = %s, want %s", classification, YahooFinancePartialMetaOnlyResponse)
		}

		data, err := mapYahooFinanceCommodityCardData(result, classification)
		if err != nil {
			t.Fatalf("mapYahooFinanceCommodityCardData() error = %v", err)
		}

		if data.Close != 81.23 || data.Price != 81.23 {
			t.Fatalf("mapped close/price = %v/%v, want 81.23", data.Close, data.Price)
		}
		if data.Open != nil {
			t.Fatalf("mapped open = %v, want nil for partial fallback", data.Open)
		}
		if data.Date != "1710259200" {
			t.Fatalf("mapped date = %q, want 1710259200", data.Date)
		}
		if data.High == nil || *data.High != 82.10 {
			t.Fatalf("mapped high = %v, want 82.10", data.High)
		}
		if data.Low == nil || *data.Low != 80.90 {
			t.Fatalf("mapped low = %v, want 80.90", data.Low)
		}
		if data.Volume == nil || *data.Volume != 456789 {
			t.Fatalf("mapped volume = %v, want 456789", data.Volume)
		}
	})

	t.Run("partial without regularMarketPrice fails boundary validation", func(t *testing.T) {
		payload := loadYahooFinanceEnvelope(t, `{
			"chart": {
				"result": [{
					"meta": {"symbol": "GC=F", "regularMarketDayHigh": 2400.0},
					"indicators": {"quote": [{}]}
				}],
				"error": null
			}
		}`)

		classification, _, err := classifyYahooFinanceResponse(payload, YahooFinanceGoldFuturesSymbol)
		if classification != YahooFinanceBoundaryFailure {
			t.Fatalf("classification = %s, want %s", classification, YahooFinanceBoundaryFailure)
		}
		if err == nil {
			t.Fatal("expected boundary failure error for missing regularMarketPrice")
		}
	})

	t.Run("symbol mismatch fails boundary validation", func(t *testing.T) {
		payload := loadYahooFinanceEnvelope(t, `{
			"chart": {
				"result": [{
					"meta": {"symbol": "SI=F"},
					"timestamp": [1710000000],
					"indicators": {"quote": [{"close": [30.0]}]}
				}],
				"error": null
			}
		}`)

		classification, _, err := classifyYahooFinanceResponse(payload, YahooFinanceGoldFuturesSymbol)
		if classification != YahooFinanceBoundaryFailure {
			t.Fatalf("classification = %s, want %s", classification, YahooFinanceBoundaryFailure)
		}
		if err == nil {
			t.Fatal("expected boundary failure error for symbol mismatch")
		}
	})
}

func TestFetchGoldPriceUsesCanonicalRestoreRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v8/finance/chart/GC=F" {
			t.Fatalf("path = %s, want /v8/finance/chart/GC=F", r.URL.Path)
		}
		if r.URL.RawQuery != yahooFinanceRestoreQuery {
			t.Fatalf("raw query = %q, want %q", r.URL.RawQuery, yahooFinanceRestoreQuery)
		}
		if cookie := r.Header.Get("Cookie"); cookie != "" {
			t.Fatalf("cookie header = %q, want empty", cookie)
		}

		fmt.Fprint(w, `{
			"chart": {
				"result": [{
					"meta": {"symbol": "GC=F"},
					"timestamp": [1710000000],
					"indicators": {
						"quote": [{
							"open": [2300.1],
							"high": [2310.0],
							"low": [2295.0],
							"close": [2308.0],
							"volume": [100]
						}]
					}
				}],
				"error": null
			}
		}`)
	}))
	defer server.Close()

	response, err := fetchGoldPriceWithClient(server.Client(), server.URL+"/v8/finance/chart/GC=F")
	if err != nil {
		t.Fatalf("fetchGoldPriceWithClient() error = %v", err)
	}
	if response.Close != 2308.0 || response.Price != 2308.0 {
		t.Fatalf("response close/price = %v/%v, want 2308", response.Close, response.Price)
	}
}

func TestFetchOilPriceUsesCanonicalRestoreRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v8/finance/chart/CL=F" {
			t.Fatalf("path = %s, want /v8/finance/chart/CL=F", r.URL.Path)
		}
		if r.URL.RawQuery != yahooFinanceRestoreQuery {
			t.Fatalf("raw query = %q, want %q", r.URL.RawQuery, yahooFinanceRestoreQuery)
		}
		if cookie := r.Header.Get("Cookie"); cookie != "" {
			t.Fatalf("cookie header = %q, want empty", cookie)
		}

		fmt.Fprint(w, `{
			"chart": {
				"result": [{
					"meta": {
						"symbol": "CL=F",
						"regularMarketPrice": 81.23,
						"regularMarketDayHigh": 82.10,
						"regularMarketDayLow": 80.90,
						"regularMarketVolume": 456789,
						"regularMarketTime": 1710259200,
						"previousClose": 80.75
					},
					"indicators": {"quote": [{}]}
				}],
				"error": null
			}
		}`)
	}))
	defer server.Close()

	response, err := fetchOilPriceWithClient(server.Client(), server.URL+"/v8/finance/chart/CL=F")
	if err != nil {
		t.Fatalf("fetchOilPriceWithClient() error = %v", err)
	}
	if response.Close != 81.23 || response.Price != 81.23 {
		t.Fatalf("response close/price = %v/%v, want 81.23", response.Close, response.Price)
	}
	if response.Open != nil {
		t.Fatalf("response open = %v, want nil for partial fallback", response.Open)
	}
}

func TestFetchGoldPriceLive(t *testing.T) {
	response, err := FetchGoldPrice()
	if err != nil {
		t.Fatalf("FetchGoldPrice() error = %v", err)
	}
	if response == nil {
		t.Fatal("FetchGoldPrice() returned nil response")
	}
	if response.Source != yahooFinanceSource {
		t.Fatalf("source = %q, want %q", response.Source, yahooFinanceSource)
	}
	if response.Price == 0 || response.Close == 0 {
		t.Fatalf("price/close = %v/%v, want non-zero", response.Price, response.Close)
	}
}

func TestFetchOilPriceLive(t *testing.T) {
	response, err := FetchOilPrice()
	if err != nil {
		t.Fatalf("FetchOilPrice() error = %v", err)
	}
	if response == nil {
		t.Fatal("FetchOilPrice() returned nil response")
	}
	if response.Source != yahooFinanceSource {
		t.Fatalf("source = %q, want %q", response.Source, yahooFinanceSource)
	}
	if response.Price == 0 || response.Close == 0 {
		t.Fatalf("price/close = %v/%v, want non-zero", response.Price, response.Close)
	}
}

func loadYahooFinanceEnvelope(t *testing.T, raw string) *yahooFinanceChartEnvelope {
	t.Helper()

	var payload yahooFinanceChartEnvelope
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	return &payload
}
