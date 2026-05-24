package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	YahooFinanceGoldFuturesSymbol = "GC=F"
	YahooFinanceOilFuturesSymbol  = "CL=F"

	YahooFinanceRestoreRange    = "5d"
	YahooFinanceRestoreInterval = "1d"

	yahooFinanceGoldChartURL = "https://query1.finance.yahoo.com/v8/finance/chart/GC=F"
	yahooFinanceOilChartURL  = "https://query1.finance.yahoo.com/v8/finance/chart/CL=F"

	yahooFinanceRestoreQuery = "range=5d&interval=1d"
	yahooFinanceSource       = "Yahoo Finance"
	yahooFinanceUserAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.178 Safari/537.36"
)

type YahooFinanceResponseClassification string

const (
	YahooFinanceFullOHLCResponse        YahooFinanceResponseClassification = "YahooFinanceFullOHLCResponse"
	YahooFinancePartialMetaOnlyResponse YahooFinanceResponseClassification = "YahooFinancePartialMetaOnlyResponse"
	YahooFinanceBoundaryFailure         YahooFinanceResponseClassification = "YahooFinanceBoundaryFailure"
)

type yahooFinanceCommodityCardData struct {
	Price  float64
	Date   string
	Source string
	Open   *float64
	Low    *float64
	Close  float64
	Volume *int64
	High   *float64
	Prior  *float64
}

type yahooFinanceChartEnvelope struct {
	Chart yahooFinanceChart `json:"chart"`
}

type yahooFinanceChart struct {
	Result []yahooFinanceChartResult `json:"result"`
	Error  *yahooFinanceChartError   `json:"error"`
}

type yahooFinanceChartError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type yahooFinanceChartResult struct {
	Meta       yahooFinanceMeta       `json:"meta"`
	Timestamp  []int64                `json:"timestamp"`
	Indicators yahooFinanceIndicators `json:"indicators"`
}

type yahooFinanceMeta struct {
	Symbol               string   `json:"symbol"`
	RegularMarketTime    *int64   `json:"regularMarketTime"`
	RegularMarketPrice   *float64 `json:"regularMarketPrice"`
	RegularMarketDayHigh *float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow  *float64 `json:"regularMarketDayLow"`
	RegularMarketVolume  *int64   `json:"regularMarketVolume"`
	PreviousClose        *float64 `json:"previousClose"`
	ChartPreviousClose   *float64 `json:"chartPreviousClose"`
}

type yahooFinanceIndicators struct {
	Quote []yahooFinanceQuote `json:"quote"`
}

type yahooFinanceQuote struct {
	Open   []*float64 `json:"open"`
	High   []*float64 `json:"high"`
	Low    []*float64 `json:"low"`
	Close  []*float64 `json:"close"`
	Volume []*int64   `json:"volume"`
}

func fetchYahooFinanceCommodityCard(client *http.Client, endpoint, expectedSymbol string) (*yahooFinanceCommodityCardData, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = yahooFinanceRestoreQuery
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", yahooFinanceUserAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newYahooFinanceBoundaryFailure("http status %s", res.Status)
	}

	var payload yahooFinanceChartEnvelope
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, newYahooFinanceBoundaryFailure("invalid json response: %v", err)
	}

	classification, result, err := classifyYahooFinanceResponse(&payload, expectedSymbol)
	if err != nil {
		return nil, err
	}

	return mapYahooFinanceCommodityCardData(result, classification)
}

func classifyYahooFinanceResponse(payload *yahooFinanceChartEnvelope, expectedSymbol string) (YahooFinanceResponseClassification, *yahooFinanceChartResult, error) {
	if payload == nil {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("missing chart payload")
	}

	if payload.Chart.Error != nil {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("chart error %s: %s", payload.Chart.Error.Code, payload.Chart.Error.Description)
	}

	if len(payload.Chart.Result) == 0 {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("missing chart result")
	}

	result := &payload.Chart.Result[0]
	if result.Meta.Symbol != expectedSymbol {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("symbol mismatch: expected %s got %s", expectedSymbol, result.Meta.Symbol)
	}

	if len(result.Indicators.Quote) == 0 {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("missing quote arrays")
	}

	quote := result.Indicators.Quote[0]
	if len(result.Timestamp) > 0 {
		if latestPresentFloatIndex(quote.Close) < 0 {
			return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("missing usable close for full OHLC payload")
		}

		return YahooFinanceFullOHLCResponse, result, nil
	}

	if hasAnyOHLCArrays(quote) {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("unexpected quote-array shape without timestamp")
	}

	if result.Meta.RegularMarketPrice == nil {
		return YahooFinanceBoundaryFailure, nil, newYahooFinanceBoundaryFailure("missing regularMarketPrice in partial response")
	}

	return YahooFinancePartialMetaOnlyResponse, result, nil
}

func mapYahooFinanceCommodityCardData(result *yahooFinanceChartResult, classification YahooFinanceResponseClassification) (*yahooFinanceCommodityCardData, error) {
	if result == nil {
		return nil, newYahooFinanceBoundaryFailure("missing result for mapping")
	}

	quote := result.Indicators.Quote[0]

	switch classification {
	case YahooFinanceFullOHLCResponse:
		idx := latestPresentFloatIndex(quote.Close)
		if idx < 0 || idx >= len(result.Timestamp) {
			return nil, newYahooFinanceBoundaryFailure("missing aligned timestamp for latest close")
		}

		closeValue, ok := floatValueAt(quote.Close, idx)
		if !ok {
			return nil, newYahooFinanceBoundaryFailure("missing close at latest usable bar")
		}

		return &yahooFinanceCommodityCardData{
			Price:  closeValue,
			Date:   fmt.Sprintf("%d", result.Timestamp[idx]),
			Source: yahooFinanceSource,
			Open:   floatPointerAt(quote.Open, idx),
			Low:    floatPointerAt(quote.Low, idx),
			Close:  closeValue,
			Volume: intPointerAt(quote.Volume, idx),
			High:   floatPointerAt(quote.High, idx),
		}, nil
	case YahooFinancePartialMetaOnlyResponse:
		price := *result.Meta.RegularMarketPrice
		date := ""
		if result.Meta.RegularMarketTime != nil {
			date = fmt.Sprintf("%d", *result.Meta.RegularMarketTime)
		}

		priorClose := result.Meta.PreviousClose
		if priorClose == nil {
			priorClose = result.Meta.ChartPreviousClose
		}

		return &yahooFinanceCommodityCardData{
			Price:  price,
			Date:   date,
			Source: yahooFinanceSource,
			Low:    result.Meta.RegularMarketDayLow,
			Close:  price,
			Volume: result.Meta.RegularMarketVolume,
			High:   result.Meta.RegularMarketDayHigh,
			Prior:  priorClose,
		}, nil
	default:
		return nil, newYahooFinanceBoundaryFailure("unsupported classification %s", classification)
	}
}

func hasAnyOHLCArrays(quote yahooFinanceQuote) bool {
	return len(quote.Open) > 0 || len(quote.High) > 0 || len(quote.Low) > 0 || len(quote.Close) > 0
}

func latestPresentFloatIndex(values []*float64) int {
	for idx := len(values) - 1; idx >= 0; idx-- {
		if values[idx] != nil {
			return idx
		}
	}

	return -1
}

func floatPointerAt(values []*float64, idx int) *float64 {
	if idx < 0 || idx >= len(values) {
		return nil
	}

	return values[idx]
}

func intPointerAt(values []*int64, idx int) *int64 {
	if idx < 0 || idx >= len(values) {
		return nil
	}

	return values[idx]
}

func floatValueAt(values []*float64, idx int) (float64, bool) {
	if idx < 0 || idx >= len(values) || values[idx] == nil {
		return 0, false
	}

	return *values[idx], true
}

func newYahooFinanceBoundaryFailure(format string, args ...any) error {
	return fmt.Errorf("%s: %s", YahooFinanceBoundaryFailure, fmt.Sprintf(format, args...))
}
