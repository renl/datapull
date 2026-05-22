# System Overview

## Purpose

Restore the broken URA pull by using the current public Property Market Information (PMI) HTML transaction search flow for private residential properties. This is a quick-restore slice, not a platform migration.

## Major Components

| Component | Responsibility | Non-Responsibilities |
|-----------|---------------|---------------------|
| URA session bootstrapper | Start a fresh PMI browser-style session, capture current `_csrf`, and retain returned cookies for the same session | Persist cookies across runs, invent tokens, bypass URA session controls |
| URA search request builder | Submit the canonical PMI form fields to the verified wire-level search URL | Use deprecated URA endpoints, synthesize unsupported fields |
| URA HTML result parser | Detect success vs no-result responses and extract transaction rows from the returned HTML fragment | Parse unrelated URA screens, depend on download flows for the restore slice |
| UI refresh path | Trigger the restore flow on startup and manual refresh, then present rows in the existing table UI | Redesign the desktop UI, change non-URA data sources |

## Component Boundaries

- The system shall treat URA as an external, brittle, session-bound HTML integration.
- The restore path shall use a two-step interaction: bootstrap `GET`, then search `POST` within the same cookie jar.
- The system shall not hard-code `_csrf`, `JSESSIONID`, `__nxquid`, or any analytics cookies.
- The system shall parse only the transaction result fragment documented in `integration-boundaries/ura-pmi-html-flow.md`.

## Decisions & Trade-offs

| Decision | Chosen Approach | Alternatives Considered | Rationale |
|----------|----------------|------------------------|-----------|
| Boundary to restore | Current public PMI HTML flow | Legacy captured flow, download endpoint, private/paid APIs | Official public flow is currently live and directly verifiable |
| Session handling | Fresh bootstrap per refresh cycle | Reuse stale cookies and tokens | The live interface is CSRF- and session-bound |
| Parse target | HTML result fragment under `#resultList` | CSV download flow | Search result HTML is enough for the current table restore |
| Scope | Minimal URA pull restore only | Broader data-platform redesign | User requested docs-only quick restore slice |

## Open Questions

| Question | Impact | Blocking? |
|----------|--------|-----------|
| Whether the implementation should clamp requests to the current 60-month window when a user-configured start date is older | Prevents silent no-data or invalid ranges | No |
| Whether tracked projects remain fixed in code or move to user-configurable state later | Affects future extensibility, not the restore contract | No |

## Verified Source Basis

- Official URA landing page for private residential property data: confirms the PMI search surface is public and current.
- Official live PMI search page at `https://eservice.ura.gov.sg/property-market-information/pmiResidentialTransactionSearch`: used to verify current form fields, hidden fields, `_csrf`, cookies, result markup, and no-result marker.
- Official URA data release calendar: used to verify update cadence.
- Optional historical reference implementation: `LearnTest` GitHub files show the same endpoint family and form pattern, but they are historical corroboration only and not authoritative for current cookies, hosts, or token values.
