# URA Refresh Flow

## Entry Point

- Application startup refresh
- User-triggered manual refresh from the existing desktop UI

## Flow

1. UI refresh path requests a URA refresh.
2. URA session bootstrapper performs `GET https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`.
3. Bootstrapper extracts current `_csrf` and retains returned session cookies in the active cookie jar.
4. Search request builder constructs the canonical form payload for each tracked project using the verified field names.
5. Client performs `POST https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch` in the same session.
6. HTML result parser classifies the response:
   - success result with `#resultList` and transaction table
   - no-result fragment with `Your search has not generated any result. Please refine your search filters.`
   - boundary failure
7. On success, parser extracts canonical result columns and emits table-ready rows.
8. UI layer replaces the URA loading placeholder with the refreshed table.
9. On no-result, UI layer shows an empty-state or no-data outcome rather than stale rows.
10. On boundary failure, UI layer shows fetch failure for URA and does not invent fallback data.

## Failure Handling

- If `_csrf` is missing, the flow shall fail boundary validation and re-bootstrap once.
- If the `POST` response lacks both success markers and the exact no-result marker, the flow shall be treated as a boundary break.
- If session state appears stale, the flow shall discard the cookie jar and retry from step 2 once.

## Sequence Summary

```text
UI -> URA session bootstrapper: refresh request
URA session bootstrapper -> URA PMI page: GET /property-market-information/pmiResidentialTransactionSearch
URA PMI page -> URA session bootstrapper: HTML page + _csrf + session cookies
URA session bootstrapper -> URA search request builder: current _csrf + cookie jar
URA search request builder -> URA PMI page: POST /property-market-information/pmiResidentialTransactionSearch
URA PMI page -> URA HTML result parser: HTML fragment
URA HTML result parser -> UI: rows | no-data | boundary failure
```

## Notes

- This flow deliberately stays on the public PMI HTML boundary.
- This flow is suitable for quick restore only and remains brittle until a more durable upstream contract is adopted.
