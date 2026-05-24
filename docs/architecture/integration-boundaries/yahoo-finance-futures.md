# Yahoo Finance Futures

## Boundary Summary

This boundary defines the restore contract for the gold and oil cards using Yahoo Finance futures chart data. It is limited to the existing provider and the existing symbols `GC=F` and `CL=F`. The system shall restore live card data through the verified chart JSON boundary only. This is a remediation slice, not a provider migration.

## Source Authority

- Official Yahoo Finance help page: `https://help.yahoo.com/kb/finance-for-web/SLN2310.html`
- Official quote page for gold futures: `https://finance.yahoo.com/quote/GC=F`
- Official quote page for crude oil futures: `https://finance.yahoo.com/quote/CL=F`
- Verified live chart endpoint for gold: `https://query1.finance.yahoo.com/v8/finance/chart/GC=F`
- Verified live chart endpoint for crude oil: `https://query1.finance.yahoo.com/v8/finance/chart/CL=F`
- Reference implementation: `https://raw.githubusercontent.com/ranaroussi/yfinance/main/yfinance/scrapers/history.py`
- Reference implementation: `https://raw.githubusercontent.com/piquette/finance-go/master/chart/client.go`

The official help page confirms Yahoo Finance covers COMEX and NYMEX market data and includes exchange-delay disclaimers. The live chart endpoint behavior is the authoritative source for this restore contract.

## Boundary Scope

In scope:

- Gold card restore for symbol `GC=F`
- Oil card restore for symbol `CL=F`
- Startup refresh and manual Refresh
- Full OHLC parsing when present
- Partial/meta-only fallback when OHLC arrays are missing

Out of scope:

- Provider migration
- Browser-page scraping from quote HTML
- Authenticated Yahoo user flows
- Historical backfill redesign
- Other Yahoo-backed cards except where this file is reused later

## Wire-Level URLs

| Purpose | Method | Wire-level URL |
|---------|--------|----------------|
| Gold futures chart data | `GET` | `https://query1.finance.yahoo.com/v8/finance/chart/GC=F` |
| Oil futures chart data | `GET` | `https://query1.finance.yahoo.com/v8/finance/chart/CL=F` |

These are the exact wire-level URLs the client connects to. No additional transport suffix or alternate mount path is involved.

## Authentication and Session Model

This restore boundary is anonymous.

- The restore path shall not require Yahoo login.
- The restore path shall not require a crumb token.
- The restore path shall not require a hard-coded Yahoo `Cookie` header.
- The restore path may send ordinary browser-like headers such as `User-Agent` and `Accept`, but those headers are not authentication.
- Any implementation that only succeeds because of a copied personal cookie jar is outside contract.

### Hard rule

Hard-coded Yahoo cookies must not be relied on for the restore path.

## Canonical Names

All names in this section are authoritative.

### Canonical symbols

- `YahooFinanceGoldFuturesSymbol` = `GC=F`
- `YahooFinanceOilFuturesSymbol` = `CL=F`

### Canonical request parameter names

- `range`
- `interval`

### Canonical allowed request parameter values for restore

- `YahooFinanceRestoreRange` = `5d`
- `YahooFinanceRestoreInterval` = `1d`

### Canonical disallowed request parameter patterns for restore

- `YahooFinanceDisallowedRestorePattern_OneDayDaily` = `range=1d&interval=1d`
- `YahooFinanceDisallowedRestorePattern_OneDayIntraday` = `range=1d&interval=1m`

### Canonical response object names

- `chart`
- `result`
- `error`
- `meta`
- `timestamp`
- `indicators`
- `quote`
- `adjclose`

### Canonical `meta` field names

- `currency`
- `symbol`
- `exchangeName`
- `fullExchangeName`
- `instrumentType`
- `regularMarketTime`
- `regularMarketPrice`
- `regularMarketDayHigh`
- `regularMarketDayLow`
- `regularMarketVolume`
- `chartPreviousClose`
- `previousClose`
- `dataGranularity`
- `range`
- `validRanges`
- `shortName`
- `priceHint`

### Canonical quote-array field names

- `open`
- `high`
- `low`
- `close`
- `volume`

### Canonical fallback field names

- `regularMarketPrice`
- `regularMarketDayHigh`
- `regularMarketDayLow`
- `regularMarketVolume`
- `regularMarketTime`
- `previousClose`
- `chartPreviousClose`

### Canonical response classifications

- `YahooFinanceFullOHLCResponse`
- `YahooFinancePartialMetaOnlyResponse`
- `YahooFinanceBoundaryFailure`

## Request Contract

### Allowed restore requests

For live card restoration, the system shall use only these requests:

| Symbol | Method | URL | Required query pattern |
|--------|--------|-----|------------------------|
| `GC=F` | `GET` | `https://query1.finance.yahoo.com/v8/finance/chart/GC=F` | `range=5d&interval=1d` |
| `CL=F` | `GET` | `https://query1.finance.yahoo.com/v8/finance/chart/CL=F` | `range=5d&interval=1d` |

### Disallowed restore requests

The following query patterns are disallowed for the restore path:

| Pattern | Status | Why disallowed |
|---------|--------|----------------|
| `range=1d&interval=1d` | Disallowed | Verified live responses for both `GC=F` and `CL=F` returned `meta` plus `indicators.quote[0] = {}` and no `timestamp` array, which is insufficient for the current card contract |
| `range=1d&interval=1m` | Disallowed | Verified live responses for both `GC=F` and `CL=F` also returned `meta` plus empty quote object and no `timestamp`, so this pattern is not contract-safe for restore |

### Required request behavior

- The client shall send `GET`.
- The client shall send the canonical query parameters exactly as `range=5d&interval=1d`.
- The client shall not assume that the bare endpoint or a `1d` range will contain usable OHLC arrays.
- The client shall treat symbol-specific URLs as distinct resources but identical contract shape.

## Response Contract

### Baseline envelope

A response is within boundary only if it is valid JSON with this top-level envelope:

- object field `chart`
- inside `chart`, field `result` or field `error`

### Success envelope

The happy-path success envelope is:

- HTTP status `200 OK`
- `chart.error = null`
- `chart.result` is a non-empty array
- first result object contains `meta`

### Full OHLC response

Classification: `YahooFinanceFullOHLCResponse`

Required conditions:

1. `chart.error = null`
2. `chart.result[0].meta` exists
3. `chart.result[0].timestamp` exists and contains at least one entry
4. `chart.result[0].indicators.quote[0]` exists
5. At least one of `open`, `high`, `low`, `close`, or `volume` contains aligned array data for the latest usable bar
6. `close` shall be present for the latest usable bar used for card rendering

Verified example behavior:

- `range=5d&interval=1d` for `GC=F` returned `timestamp`, `open`, `high`, `low`, `close`, `volume`, and `adjclose` arrays.
- `range=5d&interval=1d` for `CL=F` returned the same shape.

### Partial/meta-only response

Classification: `YahooFinancePartialMetaOnlyResponse`

Required conditions:

1. `chart.error = null`
2. `chart.result[0].meta` exists
3. `chart.result[0].indicators.quote[0]` exists but is empty or missing all OHLC arrays
4. `timestamp` is missing or empty

Verified example behavior:

- `range=1d&interval=1d` for `GC=F` returned `meta`, `indicators.quote[0] = {}`, and no `timestamp`.
- `range=1d&interval=1d` for `CL=F` returned the same partial structure.
- `range=1d&interval=1m` for both symbols also returned the same partial structure.

### Boundary failure response

Classification: `YahooFinanceBoundaryFailure`

Any of the following is boundary failure:

- non-JSON response
- HTTP non-200 response for the restore request
- `chart.error` is non-null
- `chart.result` missing or empty
- symbol mismatch between request symbol and `meta.symbol`
- JSON shape drift that removes both the quote arrays and the documented fallback `meta` fields

## Mapping Contract for Gold and Oil Cards

### Canonical latest-bar extraction rule

When classified as `YahooFinanceFullOHLCResponse`, the mapper shall use the latest index for which `close` is present. The same index position shall be used for `open`, `high`, `low`, and `volume` when present.

### Canonical fallback rule for partial responses

When classified as `YahooFinancePartialMetaOnlyResponse`, the mapper shall derive card data only from these fallback fields:

| Card field | Required fallback source order |
|-----------|--------------------------------|
| display price / close | `meta.regularMarketPrice` |
| high | `meta.regularMarketDayHigh` |
| low | `meta.regularMarketDayLow` |
| volume | `meta.regularMarketVolume` |
| timestamp/date surrogate | `meta.regularMarketTime` |
| prior close context | `meta.previousClose`, else `meta.chartPreviousClose` |

### Required behavior when fallback fields are incomplete

- If `regularMarketPrice` is missing in a partial response, the card shall fail the commodity fetch instead of inventing a close.
- Missing `regularMarketDayHigh`, `regularMarketDayLow`, or `regularMarketVolume` may render as unavailable values, but shall not overwrite `regularMarketPrice`.
- The mapper shall not substitute stale values from previous refresh cycles.

## Error Semantics

| Condition | Classification | Required handling |
|-----------|----------------|------------------|
| Full OHLC payload | `YahooFinanceFullOHLCResponse` | Render card from latest usable OHLC bar |
| Partial/meta-only payload with `regularMarketPrice` | `YahooFinancePartialMetaOnlyResponse` | Render card using documented fallback fields |
| Partial/meta-only payload without `regularMarketPrice` | `YahooFinanceBoundaryFailure` | Show fetch failure for that card |
| `chart.error` non-null | `YahooFinanceBoundaryFailure` | Show fetch failure for that card |
| HTTP non-200 | `YahooFinanceBoundaryFailure` | Show fetch failure for that card |
| Symbol mismatch | `YahooFinanceBoundaryFailure` | Reject response as invalid |

## Explicit Mismatches With Current Code

1. **Current code uses a disallowed query pattern.**
   - `api/gold.go` calls `https://query1.finance.yahoo.com/v8/finance/chart/GC=F?range=1d&interval=1d`.
   - `api/oil.go` calls `https://query1.finance.yahoo.com/v8/finance/chart/CL=F?range=1d&interval=1d`.
   - Verified live responses for that pattern are partial/meta-only and do not provide the arrays the code indexes.

2. **Current code assumes OHLC arrays are always present.**
   - Both fetchers require `Indicators.Quote[0].Close[0]` and then index `Open[0]`, `Low[0]`, `Close[0]`, `Volume[0]`, and `High[0]`.
   - Verified live partial responses return `quote[0] = {}` with no arrays.

3. **Current oil restore path relies on a hard-coded Yahoo cookie header.**
   - `api/oil.go` injects a long copied `cookie` header.
   - The restore contract forbids dependence on hard-coded Yahoo cookies.

4. **Current code does not implement partial-response fallback fields.**
   - Verified live partial responses still contain usable `meta.regularMarketPrice`, `meta.regularMarketDayHigh`, `meta.regularMarketDayLow`, `meta.regularMarketVolume`, and `meta.regularMarketTime`.
   - Current code discards these fields and returns `no data found` instead.

5. **Current response DTOs expose a `Price` field but current fetchers never populate it.**
   - `GoldPriceResponse.Price` and `OilPriceResponse.Price` exist.
   - Current fetchers return `Close` only.
   - This is not a wire-level mismatch, but it is a local contract ambiguity that downstream implementation must resolve consistently.

## Reference Implementation Notes

- `yfinance` validates chart requests against `validRanges`, distinguishes historical query modes, and treats empty quote payloads as missing-price conditions.
- `piquette/finance-go` uses the same wire-level `v8/finance/chart/<symbol>` boundary and requires quote arrays before constructing bars.

These references support the chosen endpoint family, but the live Yahoo response shapes above remain authoritative for this restore slice.

## Spike Recommendation

Smallest proof-of-contract spike:

1. Request `https://query1.finance.yahoo.com/v8/finance/chart/GC=F?range=5d&interval=1d`.
2. Confirm `chart.result[0].meta.symbol = GC=F` and quote arrays are present.
3. Request `https://query1.finance.yahoo.com/v8/finance/chart/CL=F?range=5d&interval=1d`.
4. Confirm `chart.result[0].meta.symbol = CL=F` and quote arrays are present.
5. Request the disallowed pattern `range=1d&interval=1d` for one symbol and confirm the classifier labels it `YahooFinancePartialMetaOnlyResponse`.

This spike validates the restore contract before broader commodity-card refactoring.
