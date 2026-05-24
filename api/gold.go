package api

import "net/http"

type GoldPriceResponse struct {
	Price  float64  `json:"price"`
	Date   string   `json:"date"`
	Source string   `json:"source"`
	Open   *float64 `json:"open,omitempty"`
	Low    *float64 `json:"low,omitempty"`
	Close  float64  `json:"close"`
	Volume *int64   `json:"volume,omitempty"`
	High   *float64 `json:"high,omitempty"`
}

func FetchGoldPrice() (*GoldPriceResponse, error) {
	return fetchGoldPriceWithClient(http.DefaultClient, yahooFinanceGoldChartURL)
}

func fetchGoldPriceWithClient(client *http.Client, endpoint string) (*GoldPriceResponse, error) {
	data, err := fetchYahooFinanceCommodityCard(client, endpoint, YahooFinanceGoldFuturesSymbol)
	if err != nil {
		return nil, err
	}

	return &GoldPriceResponse{
		Price:  data.Price,
		Date:   data.Date,
		Source: data.Source,
		Open:   data.Open,
		Low:    data.Low,
		Close:  data.Close,
		Volume: data.Volume,
		High:   data.High,
	}, nil
}
