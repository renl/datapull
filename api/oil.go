package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency string `json:"currency"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
					High   []float64 `json:"high"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

type OilPriceResponse struct {
	Price  float64 `json:"price"`
	Date   string  `json:"date"`
	Source string  `json:"source"`
	Open   float64 `json:"open"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
	High   float64 `json:"high"`
}

func FetchOilPrice() (*OilPriceResponse, error) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/CL=F?range=1d&interval=1d"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("cookie", "axids=gam=y-ztiAk5FE2uJ54sfgAuMZcpBktg1IVcW.~A&dv360=eS1jWm1la3Q1RTJ1SGkzOXFRR2tGdXh1ZTNaS1dyUTRJOH5B&ydsp=y-Zn85eWFE2uKEVRxOH66Wv19i6bJ05X7R~A&tbla=y-lZxDfJhE2uJBYcECmKieES.7kPRwOZUp~A; tbla_id=eb0015ff-0861-4710-b9fb-be73aa58103e-tucta9eeec0; GUC=AQEBCAFnmOtnwUIlhgUO&s=AQAAACMmCfO2&g=Z5eivg; A1=d=AQABBGdQjmECEJUoqsJa-qdVPraNUqbq9-YFEgEBCAHrmGfBZ6-0b2UB_eMBAAcIZ1COYabq9-Y&S=AQAAAkNKnjBG3dsyiFYPEqY6jjE; A3=d=AQABBGdQjmECEJUoqsJa-qdVPraNUqbq9-YFEgEBCAHrmGfBZ6-0b2UB_eMBAAcIZ1COYabq9-Y&S=AQAAAkNKnjBG3dsyiFYPEqY6jjE; A1S=d=AQABBGdQjmECEJUoqsJa-qdVPraNUqbq9-YFEgEBCAHrmGfBZ6-0b2UB_eMBAAcIZ1COYabq9-Y&S=AQAAAkNKnjBG3dsyiFYPEqY6jjE; gpp=DBAA; gpp_sid=-1; cmp=t=1739963664&j=0&u=1---; PRF=t%3DCL%253DF%252BES%253DF%252B%255EGSPC%252B%255ETYX%252BTLT%252BJPY%253DX%252BSGDUSD%253DX%252BD05.SI%252BQUBT; _chartbeat4=t=B7D9qeCid9g2OrqteBCZHfBCdLvLj&E=13&x=0&c=16.46&y=2432&w=851")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Opera GX\";v=\"114\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "Windows")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.178 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch oil price: %s", res.Status)
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

	oilPriceResponse := &OilPriceResponse{
		Date:   fmt.Sprintf("%d", date),
		Source: "Yahoo Finance",
		Open:   quote.Open[0],
		Low:    quote.Low[0],
		Close:  quote.Close[0],
		Volume: quote.Volume[0],
		High:   quote.High[0],
	}

	return oilPriceResponse, nil
}
