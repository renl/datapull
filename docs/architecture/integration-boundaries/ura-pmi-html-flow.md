# URA PMI HTML Flow

## Boundary Summary

This boundary defines the current public URA PMI HTML flow for private residential transaction search. It is brittle, HTML-coupled, and session-bound. The system shall use it only as a quick-restore path.

## Source Authority

- Official public page: `https://www.ura.gov.sg/Corporate/Property/Property-Data/Private-Residential-Properties`
- Official live PMI page and verified wire-level endpoint: `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`
- Official update cadence: `https://www.ura.gov.sg/Corporate/Property/Property-Data/Data-Release-Calendar`
- Reference implementation cited for historical pattern only: `LearnTest` URA PMI capture and download script

## Wire-Level URLs

| Purpose | Method | Wire-level URL |
|---------|--------|----------------|
| Session bootstrap and search form retrieval | `GET` | `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch` |
| Transaction search submission | `POST` | `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch` |
| Project/location popup population | `GET` | `https://eservice.ura.gov.sg/property-market-information/pmiSearchResidentialTransactionLocationPopup` |
| Result download action present in returned markup but out of scope for this restore slice | `POST` | `https://eservice.ura.gov.sg/property-market-information/pmiSearchResidentialTransactionDownload` |

## Authentication and Session Model

This is an anonymous public flow. There is no user login. The boundary is still session-bound.

### Required session bootstrap behavior

1. The client shall start with a fresh `GET https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`.
2. The client shall capture the current `_csrf` value from both:
   - `meta[name="_csrf"]`
   - `input[name="_csrf"]`
3. The client shall retain cookies set by the bootstrap response in the same cookie jar used for the subsequent `POST`.
4. The client shall treat at least these cookies as session state when present:
   - `JSESSIONID`
   - `__nxquid`
5. The client may also receive analytics cookies such as `_ga`, `_gid`, `_ga_1G0BJMEQ9S`, `_ga_KPV8HH8V5V`; these are not authoritative requirements for the business flow and must not be hard-coded.

### Hard rules

- `_csrf` must not be hard-coded.
- `JSESSIONID` must not be hard-coded.
- `__nxquid` must not be hard-coded.
- The `POST` shall reuse the same session established by the bootstrap `GET`.
- A stale token or stale cookie jar shall be treated as a recoverable boundary failure and must trigger a fresh bootstrap before retry.

## Canonical Names

### Canonical hidden/request fields

The following names are authoritative and shall be used exactly:

- `resultPerPage`
- `displayResult`
- `displayResultHeader`
- `loadAnalysis`
- `displayAnalysis`
- `displayChart`
- `displayAnalysisFilters`
- `dashboardDisplay`
- `locationDetails`
- `propertyTypeGroupNo`
- `saleYearFrom`
- `saleMonthFrom`
- `saleYearTo`
- `saleMonthTo`
- `saleType`
- `_saleType`
- `_csrf`

### Canonical result-form hidden fields observed in live result markup

- `panelNo`
- `panelId`
- `panelName`
- `transactedPriceFrom`
- `transactedPriceTo`
- `pricePerUnitAreaFrom`
- `pricePerUnitAreaTo`
- `pricePerUnitAreaUOM`
- `areaFrom`
- `areaTo`
- `areaUOM`
- `blockHouseNumber`
- `levelFrom`
- `levelTo`
- `unitNumberFrom`
- `unitNumberTo`
- `saleType[0]`
- `saleType[1]`
- `saleType[2]`
- `typeofAreaLand`
- `typeofAreaStrata`
- `enblocYes`
- `enblocNo`
- `page`
- `gotoPage`
- `tableDisplay`
- `sortBy`
- `sortAsc`
- `downloadType`
- `variableNo`
- `dataSet1No`
- `dataSet2No`
- `_selectColumn`

### Canonical response markers

- `#resultList`
- `#resultPanel1`
- `#resultForm1`
- `.downloadCSV`
- `.downloadExcel`
- `Your search has not generated any result. Please refine your search filters.`

### Canonical result column names

- `Project Name`
- `Transacted Price ($)`
- `Area (SQFT)`
- `Unit Price ($ PSF)`
- `Sale Date`
- `Street Name`
- `Type of Sale`
- `Type of Area`
- `Area (SQM)`
- `Unit Price ($ PSM)`
- `Nett Price($)`
- `Property Type`
- `Number of Units`
- `Tenure`
- `Postal District`
- `Market Segment`
- `Floor Level`

## Request Contract

### Bootstrap request

- Method: `GET`
- URL: `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`
- Purpose: obtain a fresh session and current `_csrf`

### Search request

- Method: `POST`
- URL: `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`
- Content-Type: `application/x-www-form-urlencoded; charset=UTF-8`
- Same-session cookies from bootstrap `GET`: required
- `X-Requested-With: XMLHttpRequest`: observed in live AJAX pattern and shall be sent

### Minimum form shape for restore slice

The request shall include these fields:

| Field | Required value rule |
|-------|---------------------|
| `resultPerPage` | Positive integer; live form default is `20` |
| `displayResult` | `true` |
| `displayResultHeader` | `true` |
| `loadAnalysis` | `true` |
| `displayAnalysis` | `false` |
| `displayChart` | `true` |
| `displayAnalysisFilters` | `true` |
| `dashboardDisplay` | `false` |
| `locationDetails` | JSON-encoded array string from the live flow, e.g. `["projectName","HUNDRED TREES"]` for project-based restore |
| `saleYearFrom` | Start year within URA-supported window |
| `saleMonthFrom` | Start month; live endpoint accepts `1` and `01` |
| `saleYearTo` | End year |
| `saleMonthTo` | End month; live endpoint accepts `5` and `05` |
| `saleType` | Repeated field values `1`, `2`, `3` for New Sale, Sub Sale, Resale |
| `_saleType` | `1` |
| `_csrf` | Current token from the bootstrap page |

### Optional field rules

- `propertyTypeGroupNo` is optional for `projectName` searches.
- `propertyTypeGroupNo` becomes required for postal-district searches per live page instruction.

### Parsing assumptions for `locationDetails`

- The value is not arbitrary JSON chosen by the client; it shall match the live popup selection format.
- For the quick-restore slice, project tracking shall use the verified project-name tuple format:
  - `[
    "projectName",
    "<PROJECT_NAME_1>",
    "<PROJECT_NAME_2>",
    ...
    ]`
- Maximum selected projects on the live form is five.

## Response Contract

### Success response

- HTTP status: `200 OK`
- Content-Type: `text/html; charset=utf-8`
- Body type: HTML fragment containing the results section, not a JSON document
- Required success markers for restore parsing:
  - `#resultList`
  - `#resultForm1`
  - a `table` under `#resultList`
  - table header names matching the canonical result column names above

### No-result response

- HTTP status: `200 OK`
- HTML fragment contains the exact marker:
  - `Your search has not generated any result. Please refine your search filters.`
- No result table is required in this branch.

### Download affordances present in success markup

- `.downloadCSV`
- `.downloadExcel`

These are boundary markers only for detection. They are not required to complete the restore slice.

## Parser Contract

The parser shall:

1. Determine whether the HTML is a success result, a no-result response, or an invalid boundary response.
2. On success, read rows from the first transaction table under `#resultList`.
3. Map columns by header text, not by hard-coded index alone.
4. Preserve the canonical display strings returned by URA, including commas, dashes, date text such as `Apr-26`, and tenure strings.
5. Treat missing expected markers as a boundary break.

## Update Cadence

Per the official URA Data Release Calendar and the PMI information panel:

- `Resale and Sub-Sale Private Residential Transactions`: updated twice per week, on Tuesday and Friday.
- `New Private Residential Sale Transactions`: updated weekly, on Friday.
- If the scheduled update falls on a public holiday, update moves to the following working day.

## Known Risks and Failure Modes

- HTML structure drift may break the parser.
- CSRF token rotation will break any hard-coded token.
- Session cookie rotation will break any hard-coded cookie header.
- URA may change hidden field names, default field set, or result table layout without notice.
- Search windows older than the public 60-month window may produce no results or inconsistent behavior.
- Project names are exact-string dependent; stale tracked names may return no data.

## Explicit Mismatches With Current Code

1. **Hard-coded `_csrf` is invalid by design.** Current code sends a fixed `_csrf=830dd51a-40ac-4765-b36b-72e1576e18fc`. Live interface issues a fresh token per bootstrap session.
2. **Hard-coded `Cookie` header is invalid by design.** Current code sends stale analytics cookies and no current bootstrap-derived session cookie contract.
3. **Current code skips the required bootstrap `GET`.** Live boundary requires a fresh session and current `_csrf` before `POST`.
4. **Current code assumes direct table scraping after a blind `POST`.** That only works when session state happens to align; it is not contract-safe.
5. **Current code uses zero-padded month strings by default.** Live interface accepts both padded and unpadded months, so this is not itself a break, but the contract shall not depend on padding.
6. **Current code posts to the correct wire-level URL, but without verified session bootstrap semantics.**

## Spike Recommendation

Before broader implementation slices, run the smallest proof-of-contract spike:

1. Fresh `GET` to the search page.
2. Extract `_csrf` and session cookies.
3. `POST` a single-project query such as `locationDetails=["projectName","HUNDRED TREES"]` in the same session.
4. Assert one of two exact outcomes:
   - success markers `#resultList` + results table headers, or
   - exact no-result marker string.

This spike validates the current live boundary without committing to a larger refactor.
