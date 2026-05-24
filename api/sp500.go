package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SP500Response struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
				RegularMarketVolume  int64   `json:"regularMarketVolume"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					High []float64 `json:"high"`
					Low  []float64 `json:"low"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func FetchDataSP500() ([]string, error) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/%5EGSPC?range=5d&interval=1d"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 OPR/114.0.0.0")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response SP500Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	result := response.Chart.Result[0]
	currentIndex := result.Meta.RegularMarketPrice
	regularMarketVolume := result.Meta.RegularMarketVolume
	chartPreviousClose := result.Meta.ChartPreviousClose

	volumeInMillions := float64(regularMarketVolume) / 1_000_000 // Convert to millions with decimal places

	highs := result.Indicators.Quote[0].High
	lows := result.Indicators.Quote[0].Low

	lastWeekHigh := highs[0]
	lastWeekLow := lows[0]

	for i := 1; i < len(highs); i++ {
		if highs[i] > lastWeekHigh {
			lastWeekHigh = highs[i]
		}
		if lows[i] < lastWeekLow {
			lastWeekLow = lows[i]
		}
	}

	data := []string{
		fmt.Sprintf("Current Index: %.2f", currentIndex),
		fmt.Sprintf("Previous Close: %.2f", chartPreviousClose),
		fmt.Sprintf("Volume: %.2fM", volumeInMillions),
		fmt.Sprintf("Last Week High: %.2f", lastWeekHigh),
		fmt.Sprintf("Last Week Low: %.2f", lastWeekLow),
	}

	return data, nil
}
