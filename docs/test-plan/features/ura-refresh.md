# URA Refresh Restore Test Plan

## Slice Intent

- **Founder-visible success condition:** clicking Refresh in the app shows current URA rows again without relying on stale hard-coded session artifacts.
- **Proof metric:** both the standalone URA fetch entry point and the app refresh path return non-empty current rows from a fresh session.
- **Proxy metrics:** cleaner logs, better parsing resilience, and documentation updates.

## Integration Points

- `GET https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`
- `POST https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`
- Existing app refresh path in `ui/ui.go` and `ui/ura_table.go`
- Existing standalone URA entry point in `test_ura.go`

## Scenarios

### 1. Fresh-session bootstrap succeeds
- Start with an empty cookie jar.
- Perform the bootstrap `GET`.
- Extract current `_csrf` from the bootstrap page and retain session cookies.
- **Expected:** `_csrf` is present, the cookie jar contains current session cookies, and no hard-coded cookie/token path is required.

### 2. Search request reuses the same session
- Reuse the bootstrap client/session for the `POST`.
- Submit canonical fields for the tracked project-name query.
- **Expected:** request uses the same session, sends `X-Requested-With: XMLHttpRequest`, and receives either a success HTML fragment or the exact no-result marker.

### 3. Success response parses current rows
- Execute the restore query for the currently tracked projects.
- Parse rows from the first transaction table under `#resultList`.
- **Expected:** header mapping follows canonical column names and the result contains at least one non-header row for the current query window.

### 4. Boundary failure triggers recovery
- Simulate or detect stale/missing `_csrf` or session state.
- **Expected:** implementation performs one fresh bootstrap retry and fails clearly if boundary markers are still invalid.

### 5. No-result response is classified correctly
- Submit a query that returns no data, or validate the parsing branch with a fixture if live no-result is not convenient.
- **Expected:** implementation recognizes the exact no-result marker and does not display stale rows.

### 6. App refresh path shows fresh rows
- Launch the app and trigger Refresh.
- **Expected:** the URA section replaces the loading placeholder with refreshed rows from a fresh session, with no dependence on stale hard-coded artifacts.

## Validation Commands

```bash
gofmt -w api/ura.go test_ura.go ui/ui.go ui/ura_table.go
gofmt -l .
go test ./...
go build ./...
go run ./test_ura.go
go run .
```

## Approval Gate

- Do not approve this slice unless at least one real-endpoint URA request succeeds against the live PMI boundary.
- Mock-only evidence is insufficient.
- Proof metric is unmet if either the standalone path or the app refresh path remains unverified or returns only headers / empty rows.
