package api

import "net/http"

type OilPriceResponse struct {
	Price  float64  `json:"price"`
	Date   string   `json:"date"`
	Source string   `json:"source"`
	Open   *float64 `json:"open,omitempty"`
	Low    *float64 `json:"low,omitempty"`
	Close  float64  `json:"close"`
	Volume *int64   `json:"volume,omitempty"`
	High   *float64 `json:"high,omitempty"`
}

func FetchOilPrice() (*OilPriceResponse, error) {
	return fetchOilPriceWithClient(http.DefaultClient, yahooFinanceOilChartURL)
}

func fetchOilPriceWithClient(client *http.Client, endpoint string) (*OilPriceResponse, error) {
	data, err := fetchYahooFinanceCommodityCard(client, endpoint, YahooFinanceOilFuturesSymbol)
	if err != nil {
		return nil, err
	}

	return &OilPriceResponse{
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
