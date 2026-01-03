package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BitcoinPriceResponse struct {
	Price  float64 `json:"price"`
	Date   string  `json:"date"`
	Source string  `json:"source"`
	Open   float64 `json:"open"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
	High   float64 `json:"high"`
}

func FetchBitcoinPrice() (*BitcoinPriceResponse, error) {
	// BTC-USD is the Bitcoin to USD symbol on Yahoo Finance
	url := "https://query1.finance.yahoo.com/v8/finance/chart/BTC-USD?range=1d&interval=1d"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.178 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch bitcoin price: %s", res.Status)
	}

	var yahooResponse YahooFinanceResponse
	if err := json.NewDecoder(res.Body).Decode(&yahooResponse); err != nil {
		return nil, err
	}

	if len(yahooResponse.Chart.Result) == 0 || len(yahooResponse.Chart.Result[0].Indicators.Quote) == 0 || len(yahooResponse.Chart.Result[0].Indicators.Quote[0].Close) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	quote := yahooResponse.Chart.Result[0].Indicators.Quote[0]
	date := yahooResponse.Chart.Result[0].Timestamp[0]

	bitcoinPriceResponse := &BitcoinPriceResponse{
		Date:   fmt.Sprintf("%d", date),
		Source: "Yahoo Finance",
		Open:   quote.Open[0],
		Low:    quote.Low[0],
		Close:  quote.Close[0],
		Volume: quote.Volume[0],
		High:   quote.High[0],
	}

	return bitcoinPriceResponse, nil
}
