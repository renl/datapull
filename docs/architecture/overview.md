# System Overview

## Purpose

Restore broken live data retrieval on existing cards by using verified public upstream interfaces. Current in-scope remediation covers Yahoo Finance futures chart retrieval for gold and oil so the cards show live values again on startup and manual Refresh. This is an architecture remediation slice, not a provider migration.

## Major Components

| Component | Responsibility | Non-Responsibilities |
|-----------|---------------|---------------------|
| URA session bootstrapper | Start a fresh PMI browser-style session, capture current `_csrf`, and retain returned cookies for the same session | Persist cookies across runs, invent tokens, bypass URA session controls |
| URA search request builder | Submit the canonical PMI form fields to the verified wire-level search URL | Use deprecated URA endpoints, synthesize unsupported fields |
| URA HTML result parser | Detect success vs no-result responses and extract transaction rows from the returned HTML fragment | Parse unrelated URA screens, depend on download flows for the restore slice |
| Yahoo futures request builder | Build the allowed Yahoo Finance chart request for `GC=F` and `CL=F` using canonical query parameters | Use disallowed query patterns, depend on browser-page scraping, require hard-coded cookies |
| Yahoo futures response classifier | Classify Yahoo chart responses as full OHLC, partial/meta-only, or boundary failure | Assume quote arrays always exist, treat partial responses as valid card payloads |
| Commodity card mapper | Map canonical Yahoo response fields into the existing gold/oil card display contract with fallback precedence | Invent prices, silently substitute stale values, redesign card layout |
| UI refresh path | Trigger the restore flow on startup and manual refresh, then present rows in the existing table UI | Redesign the desktop UI, change non-URA data sources |

## Component Boundaries

- The system shall treat URA as an external, brittle, session-bound HTML integration.
- The restore path shall use a two-step interaction: bootstrap `GET`, then search `POST` within the same cookie jar.
- The system shall not hard-code `_csrf`, `JSESSIONID`, `__nxquid`, or any analytics cookies.
- The system shall parse only the transaction result fragment documented in `integration-boundaries/ura-pmi-html-flow.md`.
- The system shall treat Yahoo Finance futures chart data as an external, undocumented-but-live JSON boundary that must be verified from current live responses before use.
- The Yahoo restore path shall use the wire-level chart endpoint `https://query1.finance.yahoo.com/v8/finance/chart/<symbol>`.
- The Yahoo restore path shall not rely on any hard-coded Yahoo `Cookie` header.
- The Yahoo restore path shall distinguish full OHLC responses from partial/meta-only responses and shall apply only the documented fallback fields for partial responses.
- Startup refresh and manual Refresh shall use the same commodity-card boundary contract.

## Decisions & Trade-offs

| Decision | Chosen Approach | Alternatives Considered | Rationale |
|----------|----------------|------------------------|-----------|
| Boundary to restore | Current public PMI HTML flow | Legacy captured flow, download endpoint, private/paid APIs | Official public flow is currently live and directly verifiable |
| Session handling | Fresh bootstrap per refresh cycle | Reuse stale cookies and tokens | The live interface is CSRF- and session-bound |
| Parse target | HTML result fragment under `#resultList` | CSV download flow | Search result HTML is enough for the current table restore |
| Scope | Minimal URA pull restore only | Broader data-platform redesign | User requested docs-only quick restore slice |
| Yahoo futures restore boundary | `v8/finance/chart/<symbol>` JSON endpoint with canonical parameter pattern | Quote-page HTML scraping, provider migration, stale captured requests | Live JSON response is directly verifiable and already aligns with repo intent |
| Allowed Yahoo query pattern for restore | Daily restore request using `range=5d&interval=1d` | `range=1d&interval=1d`, `range=1d&interval=1m`, cookie-bound browser mimicry | Verified live responses show `1d` requests can return meta-only payloads while `5d/1d` returns usable OHLC arrays |
| Partial Yahoo response handling | Explicit classification plus fallback to canonical `meta` fields | Blind array indexing | Live responses can omit `timestamp` and OHLC arrays for `1d` range requests |
| Cookie strategy for Yahoo | Anonymous stateless HTTP requests without hard-coded cookies | Hard-coded browser cookies copied from a desktop session | Hard-coded cookies are unstable, non-portable, and not required for the verified restore query pattern |

## Open Questions

| Question | Impact | Blocking? |
|----------|--------|-----------|
| Whether the implementation should clamp requests to the current 60-month window when a user-configured start date is older | Prevents silent no-data or invalid ranges | No |
| Whether tracked projects remain fixed in code or move to user-configurable state later | Affects future extensibility, not the restore contract | No |
| Whether future work should introduce a shared commodity quote adapter for gold, oil, S&P 500, and bitcoin | Affects later consolidation only | No |

## Verified Source Basis

- Official URA landing page for private residential property data: confirms the PMI search surface is public and current.
- Official live PMI search page at `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`: used to verify current form fields, hidden fields, `_csrf`, cookies, result markup, and no-result marker.
- Official URA data release calendar: used to verify update cadence.
- Optional historical reference implementation: `LearnTest` GitHub files show the same endpoint family and form pattern, but they are historical corroboration only and not authoritative for current cookies, hosts, or token values.
- Official Yahoo Finance coverage/help page at `https://help.yahoo.com/kb/finance-for-web/SLN2310.html`: used to verify that Yahoo Finance covers COMEX and NYMEX market data and documents provider/delay disclaimers.
- Official Yahoo Finance quote pages at `https://finance.yahoo.com/quote/GC=F` and `https://finance.yahoo.com/quote/CL=F`: used to verify that the symbols resolve to Gold and Crude Oil futures surfaces on Yahoo Finance.
- Verified live Yahoo chart endpoint responses at `https://query1.finance.yahoo.com/v8/finance/chart/GC=F` and `https://query1.finance.yahoo.com/v8/finance/chart/CL=F`, including explicit checks of `range=1d&interval=1d`, `range=1d&interval=1m`, and `range=5d&interval=1d` behavior.
- Reference implementations: `yfinance` history scraper and `piquette/finance-go` chart client. These are corroborating references only; live endpoint behavior is authoritative for this restore slice.

## Verified Source Basis

- Official URA landing page for private residential property data: confirms the PMI search surface is public and current.
- Official live PMI search page at `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`: used to verify current form fields, hidden fields, `_csrf`, cookies, result markup, and no-result marker.
- Official URA data release calendar: used to verify update cadence.
- Optional historical reference implementation: `LearnTest` GitHub files show the same endpoint family and form pattern, but they are historical corroboration only and not authoritative for current cookies, hosts, or token values.
