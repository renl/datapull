# Gold and Oil Card Restore Test Plan

## Slice Intent

- **Founder-visible success condition:** gold and oil cards show live values again both on startup and when the user presses Refresh.
- **Proof metric:** both cards render data sourced from current live Yahoo Finance responses for `GC=F` and `CL=F` on startup and Refresh.
- **Proxy metrics:** safer Yahoo query pattern, no hard-coded Yahoo cookie dependency, and resilient parsing/classification of partial/meta-only Yahoo responses.

## Baseline and Scope

- **Slice:** `gold-oil-restore-03-integration`
- **Prerequisite:** none
- **Baseline:** `main`
- **External boundary:** `docs/architecture/integration-boundaries/yahoo-finance-futures.md`
- **Key flow:** `docs/architecture/key-flows/commodity-card-refresh.md`

## Integration Points

- `GET https://query1.finance.yahoo.com/v8/finance/chart/GC=F?range=5d&interval=1d`
- `GET https://query1.finance.yahoo.com/v8/finance/chart/CL=F?range=5d&interval=1d`
- Existing commodity fetchers in `api/gold.go` and `api/oil.go`
- Existing startup and Refresh path in `ui/ui.go`, `ui/gold_card.go`, and `ui/oil_card.go`

## Scenarios

### 1. Gold restore uses the canonical Yahoo request
- Execute the gold fetch path.
- **Expected:** the request uses `GET https://query1.finance.yahoo.com/v8/finance/chart/GC=F?range=5d&interval=1d`.
- **Expected:** no hard-coded Yahoo cookie is required.

### 2. Oil restore uses the canonical Yahoo request
- Execute the oil fetch path.
- **Expected:** the request uses `GET https://query1.finance.yahoo.com/v8/finance/chart/CL=F?range=5d&interval=1d`.
- **Expected:** no hard-coded Yahoo cookie is required.

### 3. Full OHLC response maps the latest usable bar
- Run the fetch path against a live `5d/1d` Yahoo response.
- **Expected:** classification is `YahooFinanceFullOHLCResponse` when quote arrays are present.
- **Expected:** the mapper uses the latest index where `close` is present and aligns `open`, `high`, `low`, and `volume` to that same index.

### 4. Partial/meta-only response follows documented fallback precedence
- Exercise the classifier with a live or fixture-backed `range=1d&interval=1d` response shape.
- **Expected:** classification is `YahooFinancePartialMetaOnlyResponse` when `meta` exists but quote arrays and `timestamp` are absent.
- **Expected:** fallback values use only `meta.regularMarketPrice`, `meta.regularMarketDayHigh`, `meta.regularMarketDayLow`, `meta.regularMarketVolume`, `meta.regularMarketTime`, and `meta.previousClose` else `meta.chartPreviousClose`.

### 5. Partial response without `regularMarketPrice` fails cleanly
- Exercise a partial/meta-only payload that omits `regularMarketPrice`.
- **Expected:** classification escalates to `YahooFinanceBoundaryFailure` for card rendering.
- **Expected:** the card shows fetch failure rather than inventing or reusing stale values.

### 6. Symbol mismatch or malformed envelope fails boundary validation
- Exercise a payload with missing `chart.result`, non-null `chart.error`, or a mismatched `meta.symbol`.
- **Expected:** the fetch path rejects the response as `YahooFinanceBoundaryFailure`.

### 7. Startup path restores cards independently
- Launch the app from a clean run.
- **Expected:** gold and oil placeholders are replaced independently using live Yahoo responses.
- **Expected:** failure of one commodity does not block the other card.

### 8. Manual Refresh re-runs the same commodity flow
- Launch the app and trigger Refresh.
- **Expected:** startup and Refresh use the same commodity restore contract.
- **Expected:** Refresh does not depend on stale in-memory commodity payloads from a previous cycle.

## Validation Commands

```bash
gofmt -w api/gold.go api/oil.go ui/gold_card.go ui/oil_card.go ui/ui.go
gofmt -l .
go test ./... -count=1
go build ./...
go test ./... -run 'TestFetch(Gold|Oil)PriceLive|TestYahooFinanceClassifier' -count=1
go run .
```

## Approval Gate

- Do not approve this slice unless at least one real-endpoint request succeeds for both `GC=F` and `CL=F` against the live Yahoo chart boundary.
- Mock-only evidence is insufficient because this slice changes an external system interaction.
- Proof metric remains unmet if verification shows only safer parsing/query logic without confirming both cards render live data on startup and Refresh.
