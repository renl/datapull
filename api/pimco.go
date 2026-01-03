package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PimcoResponse struct {
	Elements []struct {
		Type           string `json:"Type"`
		CompanyName    string `json:"CompanyName"`
		ComponentSeries []struct {
			Type string `json:"Type"`
			Values []float64 `json:"Values"`
		} `json:"ComponentSeries"`
	} `json:"Elements"`
}

func FetchDataPimcoGISIncome() ([]string, error) {
	url := "https://markets.ft.com/data/chartapi/series"
	method := "POST"

	payload := bytes.NewBufferString(`{"days":180,"dataNormalized":false,"dataPeriod":"Day","dataInterval":1,"realtime":false,"yFormat":"0.###","timeServiceFormat":"JSON","rulerIntradayStart":26,"rulerIntradayStop":3,"rulerInterdayStart":10957,"rulerInterdayStop":365,"returnDateType":"ISO8601","elements":[{"Label":"09f5b483","Type":"price","Symbol":"56954688","OverlayIndicators":[],"Params":{}}]}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", "__RequestVerificationToken_L2RhdGE1=xczAhKSitbT0UC40vit8ieRBq97Y0a7BEkdoOkTMWEKcKMCkvb0tlEKcWIMNB8DDJm6nfLyEKXQAaoYukgHU8W34M7OiDWqKygLdLAX7bwI1; GZIP=1; consentUUID=bdb72d0b-b7be-4fc2-869e-cadc34ed1c77_34; consentDate=2024-08-03T01:06:03.387Z; FTConsent=marketingBypost%3Aoff%2CmarketingByemail%3Aoff%2CmarketingByphonecall%3Aoff%2CmarketingByfax%3Aoff%2CmarketingBysms%3Aoff%2CenhancementBypost%3Aoff%2CenhancementByemail%3Aoff%2CenhancementByphonecall%3Aoff%2CenhancementByfax%3Aoff%2CenhancementBysms%3Aoff%2CbehaviouraladsOnsite%3Aon%2CdemographicadsOnsite%3Aon%2CrecommendedcontentOnsite%3Aon%2CprogrammaticadsOnsite%3Aon%2CcookiesUseraccept%3Aoff%2CcookiesOnsite%3Aoff%2CmembergetmemberByemail%3Aoff%2CpermutiveadsOnsite%3Aon%2CpersonalisedmarketingOnsite%3Aon; FTCookieConsentGDPR=true; spoor-id=clzdfn33900003j7cyefbypun; permutive-id=bbc81ea5-9b9b-43b8-831c-e72ce611c967; usnatUUID=cae320b7-c997-4005-bef8-054b881731c6; _gid=GA1.2.448820738.1737382142; __gads=ID=863ba6cfa2da5c0c:T=1722647163:RT=1737427144:S=ALNI_Mb20Uz3P3cXf3JDxjrccD4UrG22zA; __eoi=ID=52f79613d1bee95d:T=1722647163:RT=1737427144:S=AA-Afjaa1qzPjGv5f9V8R762t3TV; _gat_UA-165575472-1=1; _ga_2DSMN2JH8F=GS1.1.1737427141.261.1.1737427628.0.0.0; _ga=GA1.1.1270686200.1722647158")
	req.Header.Add("Origin", "https://markets.ft.com")
	req.Header.Add("Referer", "https://markets.ft.com/data/funds/tearsheet/charts?s=IE00B9HH6X13:SGD")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 OPR/114.0.0.0")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Opera GX\";v=\"114\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response PimcoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var data []string

	found := false
	for _, element := range response.Elements {
		if element.Type == "price" {
			data = append(data, element.CompanyName)
			for _, component := range element.ComponentSeries {
				if component.Type == "Open" {
					values := component.Values
					if len(values) >= 2 {
						data = append(data, fmt.Sprintf("SGD %.2f", values[len(values)-2]))
						data = append(data, fmt.Sprintf("SGD %.2f", values[len(values)-1]))
						found = true
						break
					}
				}
			}
			if (found) {
				break
			}
		}
	}

	return data, nil
}
