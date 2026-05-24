# Commodity Card Refresh Flow

## Entry Point

- Application startup refresh
- User-triggered manual Refresh from the existing desktop UI

## Scope

This flow covers only the Yahoo Finance restore path for:

- gold card using `GC=F`
- oil card using `CL=F`

## Flow

1. UI refresh path starts a commodity-card refresh for gold and oil.
2. Commodity refresh coordinator creates one restore request per symbol using the canonical wire-level URL.
3. For each symbol, Yahoo futures request builder issues `GET https://query1.finance.yahoo.com/v8/finance/chart/<symbol>?range=5d&interval=1d`.
4. Yahoo futures response classifier inspects the JSON envelope.
5. If the response matches `YahooFinanceFullOHLCResponse`, the commodity card mapper selects the latest usable `close` bar and aligned `open`, `high`, `low`, and `volume` values.
6. If the response matches `YahooFinancePartialMetaOnlyResponse`, the commodity card mapper applies fallback precedence from `meta.regularMarketPrice`, `meta.regularMarketDayHigh`, `meta.regularMarketDayLow`, `meta.regularMarketVolume`, and `meta.regularMarketTime`.
7. If the response matches `YahooFinanceBoundaryFailure`, the affected card enters fetch-failure presentation.
8. UI layer replaces each commodity loading placeholder independently so one card can succeed while the other fails.
9. On manual Refresh, the same flow repeats with no reuse of stale commodity payloads.

## Failure Handling

- If the request uses a disallowed query pattern such as `range=1d&interval=1d`, the response shall not be treated as a full restore success.
- If quote arrays are absent but `regularMarketPrice` exists, the card may render from fallback fields only.
- If both quote arrays and `regularMarketPrice` are absent, the card shall show fetch failure.
- Hard-coded Yahoo cookies shall not be added as a recovery step.
- Failure of gold retrieval shall not block oil rendering, and failure of oil retrieval shall not block gold rendering.

## Sequence Summary

```text
UI -> Commodity refresh coordinator: startup or Refresh request
Commodity refresh coordinator -> Yahoo futures request builder: build request for GC=F and CL=F
Yahoo futures request builder -> Yahoo Finance chart endpoint: GET /v8/finance/chart/<symbol>?range=5d&interval=1d
Yahoo Finance chart endpoint -> Yahoo futures response classifier: JSON response
Yahoo futures response classifier -> Commodity card mapper: FullOHLC | PartialMetaOnly | BoundaryFailure
Commodity card mapper -> UI: gold card payload | oil card payload | fetch failure per card
```

## Notes

- This flow keeps the current provider and current symbols.
- This flow exists to restore founder-visible live card data safely.
- This flow requires resilient parsing because verified live Yahoo responses can omit OHLC arrays for narrower query patterns.
