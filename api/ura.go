package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	uraSearchURL   = "https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch"
	noResultMarker = "Your search has not generated any result. Please refine your search filters."
)

var (
	errURABoundaryFailure = errors.New("ura PMI boundary failure")
	errURANoResults       = errors.New("ura PMI returned no results")

	trackedProjects = []string{"HUNDRED TREES", "SUITES @ EASTCOAST"}

	canonicalURAColumns = []string{
		"Project Name",
		"Transacted Price ($)",
		"Area (SQFT)",
		"Unit Price ($ PSF)",
		"Sale Date",
		"Street Name",
		"Type of Sale",
		"Type of Area",
		"Area (SQM)",
		"Unit Price ($ PSM)",
		"Nett Price($)",
		"Property Type",
		"Number of Units",
		"Tenure",
		"Postal District",
		"Market Segment",
		"Floor Level",
	}
)

func FetchDataURA() ([][]string, error) {
	for attempt := 0; attempt < 2; attempt++ {
		data, err := fetchDataURAFreshSession()
		if err == nil {
			return data, nil
		}

		if errors.Is(err, errURANoResults) {
			return [][]string{canonicalURAColumns}, nil
		}

		if attempt == 0 && errors.Is(err, errURABoundaryFailure) {
			continue
		}

		return nil, err
	}

	return nil, errURABoundaryFailure
}

func fetchDataURAFreshSession() ([][]string, error) {
	client, err := newURAClient()
	if err != nil {
		return nil, err
	}

	csrfToken, err := bootstrapURASession(client)
	if err != nil {
		return nil, err
	}

	body, err := searchURA(client, csrfToken)
	if err != nil {
		return nil, err
	}

	return parseURAResult(body)
}

func newURAClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	return &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}, nil
}

func bootstrapURASession(client *http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uraSearchURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", uraSearchURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("bootstrap GET failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bootstrap GET returned %s: %w", res.Status, errURABoundaryFailure)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	metaCSRF, metaExists := doc.Find(`meta[name="_csrf"]`).Attr("content")
	inputCSRF, inputExists := doc.Find(`input[name="_csrf"]`).Attr("value")

	metaCSRF = strings.TrimSpace(metaCSRF)
	inputCSRF = strings.TrimSpace(inputCSRF)

	if !metaExists || !inputExists || metaCSRF == "" || inputCSRF == "" {
		return "", fmt.Errorf("bootstrap GET missing _csrf markers: %w", errURABoundaryFailure)
	}

	if metaCSRF != inputCSRF {
		return "", fmt.Errorf("bootstrap GET returned mismatched _csrf markers: %w", errURABoundaryFailure)
	}

	bootstrapURL, err := url.Parse(uraSearchURL)
	if err != nil {
		return "", err
	}

	if len(client.Jar.Cookies(bootstrapURL)) == 0 {
		return "", fmt.Errorf("bootstrap GET returned no session cookies: %w", errURABoundaryFailure)
	}

	return inputCSRF, nil
}

func searchURA(client *http.Client, csrfToken string) ([]byte, error) {
	now := time.Now()
	windowStart := now.AddDate(0, -59, 0)

	locationDetails, err := json.Marshal(append([]string{"projectName"}, trackedProjects...))
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("resultPerPage", "50")
	form.Set("displayResult", "true")
	form.Set("displayResultHeader", "true")
	form.Set("loadAnalysis", "true")
	form.Set("displayAnalysis", "false")
	form.Set("displayChart", "true")
	form.Set("displayAnalysisFilters", "true")
	form.Set("dashboardDisplay", "false")
	form.Set("locationDetails", string(locationDetails))
	form.Set("saleYearFrom", strconv.Itoa(windowStart.Year()))
	form.Set("saleMonthFrom", strconv.Itoa(int(windowStart.Month())))
	form.Set("saleYearTo", strconv.Itoa(now.Year()))
	form.Set("saleMonthTo", strconv.Itoa(int(now.Month())))
	form.Add("saleType", "1")
	form.Add("saleType", "2")
	form.Add("saleType", "3")
	form.Set("_saleType", "1")
	form.Set("_csrf", csrfToken)

	req, err := http.NewRequest(http.MethodPost, uraSearchURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://eservice.ura.gov.sg")
	req.Header.Set("Referer", uraSearchURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search POST failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search POST returned %s: %w", res.Status, errURABoundaryFailure)
	}

	return body, nil
}

func parseURAResult(body []byte) ([][]string, error) {
	bodyText := string(body)
	if strings.Contains(bodyText, noResultMarker) {
		return nil, errURANoResults
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resultList := doc.Find("#resultList")
	resultForm := doc.Find("#resultForm1")
	table := resultList.Find("table").First()

	if resultList.Length() == 0 || resultForm.Length() == 0 || table.Length() == 0 {
		return nil, fmt.Errorf("search POST response missing success markers: %w", errURABoundaryFailure)
	}

	headers, err := parseURATableHeaders(table)
	if err != nil {
		return nil, err
	}

	dataRows, err := parseURATableRows(table, headers)
	if err != nil {
		return nil, err
	}
	if len(dataRows) == 0 {
		return nil, fmt.Errorf("result table contained no transaction rows: %w", errURABoundaryFailure)
	}

	data := make([][]string, 0, len(dataRows)+1)
	data = append(data, canonicalURAColumns)
	data = append(data, dataRows...)

	return data, nil
}

func parseURATableHeaders(table *goquery.Selection) ([]string, error) {
	headerRow := table.Find("thead tr").First()
	if headerRow.Length() == 0 {
		headerRow = table.Find("tr").First()
	}

	if headerRow.Length() == 0 {
		return nil, fmt.Errorf("result table missing header row: %w", errURABoundaryFailure)
	}

	headers := make([]string, 0)
	headerRow.Find("th,td").Each(func(i int, cell *goquery.Selection) {
		header := normalizeURAText(cell.Text())
		if header != "" {
			headers = append(headers, header)
		}
	})

	if len(headers) == 0 {
		return nil, fmt.Errorf("result table missing headers: %w", errURABoundaryFailure)
	}

	canonicalSet := make(map[string]struct{}, len(canonicalURAColumns))
	for _, column := range canonicalURAColumns {
		canonicalSet[column] = struct{}{}
	}

	for _, header := range headers {
		if _, ok := canonicalSet[header]; !ok {
			return nil, fmt.Errorf("unexpected result header %q: %w", header, errURABoundaryFailure)
		}
	}

	return headers, nil
}

func parseURATableRows(table *goquery.Selection, headers []string) ([][]string, error) {
	rows := make([][]string, 0)

	bodyRows := table.Find("tbody tr")
	if bodyRows.Length() == 0 {
		bodyRows = table.Find("tr").Slice(1, goquery.ToEnd)
	}

	bodyRows.Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 0 {
			return
		}

		mapped := make(map[string]string, len(headers))
		cells.Each(func(j int, cell *goquery.Selection) {
			if j >= len(headers) {
				return
			}
			mapped[headers[j]] = normalizeURAText(cell.Text())
		})

		rowData := make([]string, 0, len(canonicalURAColumns))
		hasValue := false
		for _, column := range canonicalURAColumns {
			value := mapped[column]
			if value != "" {
				hasValue = true
			}
			rowData = append(rowData, value)
		}

		if hasValue {
			rows = append(rows, rowData)
		}
	})

	return rows, nil
}

func normalizeURAText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}
