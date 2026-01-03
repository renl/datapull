package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func FetchDataURA() ([][]string, error) {
	url := "https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch"
	method := "POST"
	now := time.Now()
	currentYear := now.Year()
	currentMonth := fmt.Sprintf("%02d", int(now.Month()))
	payload := bytes.NewBufferString(fmt.Sprintf("resultPerPage=50&displayResult=true&displayResultHeader=true&loadAnalysis=true&displayAnalysis=false&displayChart=true&displayAnalysisFilters=true&dashboardDisplay=false&locationDetails=%%5B%%22projectName%%22%%2C%%22HUNDRED+TREES%%22%%2C%%22SUITES+%%40+EASTCOAST%%22%%5D&saleYearFrom=2020&saleMonthFrom=01&saleYearTo=%d&saleMonthTo=%s&saleType=1&saleType=2&saleType=3&_saleType=1&_csrf=830dd51a-40ac-4765-b36b-72e1576e18fc", currentYear, currentMonth))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", "_ga_YK05YF0QS2=GS1.1.1737033469.6.1.1737033511.18.0.2069279388; _ga=GA1.3.1286357919.1674983656; _sp_id.b186=79f96974-4a57-4a11-8c9a-f90821d875ec.1686051215.4.1737270975.1712323839.6b6ad702-05a2-477a-89c9-612e3cf6e79f.3ffabffe-5fb3-44bf-a8a5-8e1356411f7e.8d1c3314-9ec5-4ec8-9647-d07dd0d346a3.1737270974507.1; _gid=GA1.3.1158824259.1738773593; _gat=1; _ga_1G0BJMEQ9S=GS1.3.1738773593.167.1.1738774420.0.0.0; _ga_KPV8HH8V5V=GS1.3.1738773593.167.1.1738774420.0.0.0")
	req.Header.Add("Origin", "https://eservice.ura.gov.sg")
	req.Header.Add("Referer", "https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 OPR/114.0.0.0")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
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

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var data [][]string

	// Add headers
	headers := []string{"Project Name", "Transacted Price ($)", "Area (SQFT)", "Unit Price ($ PSF)", "Sale Date", "Street Name", "Type of Sale", "Type of Area", "Area (SQM)", "Unit Price ($ PSM)", "Nett Price($)", "Property Type", "Number of Units", "Tenure", "Postal District", "Market Segment", "Floor Level"}
	data = append(data, headers)

	// Scrape data
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			rowData = append(rowData, text)
		})
		data = append(data, rowData)
	})

	return data, nil
}
