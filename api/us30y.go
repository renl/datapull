package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type US30YResponse struct {
	FormattedQuoteResult struct {
		FormattedQuote []struct {
			Name string `json:"name"`
			Last string `json:"last"`
			Open string `json:"open"`
			High string `json:"high"`
			Low  string `json:"low"`
		} `json:"FormattedQuote"`
	} `json:"FormattedQuoteResult"`
}

func FetchDataUS30Y() ([][]string, error) {
	url := "https://quote.cnbc.com/quote-html-webservice/restQuote/symbolType/symbol?symbols=US30Y&requestMethod=itv&noform=1&partnerId=2&fund=1&exthrs=1&output=json&events=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("origin", "https://www.cnbc.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://www.cnbc.com/quotes/US30Y")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Opera GX\";v=\"114\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 OPR/114.0.0.0")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response US30YResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var data [][]string
	headers := []string{"Name", "Last", "Open", "High", "Low"}
	data = append(data, headers)

	for _, quote := range response.FormattedQuoteResult.FormattedQuote {
		row := []string{quote.Name, quote.Last, quote.Open, quote.High, quote.Low}
		data = append(data, row)
	}

	return data, nil
}
