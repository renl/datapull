# Architecture Contract

> The steering wheel for the project. Language-agnostic system design that preserves the high-level picture and makes the project easy to steer.

**Owned by:** `@software-architect`

## Structure

| File | Contents |
|------|----------|
| `overview.md` | System overview — purpose, major components, boundaries, decisions & trade-offs, open questions |
| `integration-boundaries/ura-pmi-html-flow.md` | URA public PMI HTML integration boundary for private residential transaction search |
| `integration-boundaries/yahoo-finance-futures.md` | Yahoo Finance futures chart restore boundary for gold (`GC=F`) and oil (`CL=F`) cards |
| `key-flows/ura-refresh.md` | Quick-restore refresh flow for URA transaction retrieval |
| `key-flows/commodity-card-refresh.md` | Startup and manual refresh flow for Yahoo Finance gold and oil commodity cards |

## Conventions

- Read this README first.
- Canonical names defined in integration-boundary files are authoritative.
- Transport and request URLs must use wire-level URLs.
- For each external system, read only the relevant boundary file instead of inferring behavior from unrelated integrations.

## How to use this index

1. Read this README first to orient.
2. Read `overview.md` for current architectural steering decisions.
3. For URA work, read `integration-boundaries/ura-pmi-html-flow.md` and `key-flows/ura-refresh.md`.
4. For gold/oil card restore work, read `integration-boundaries/yahoo-finance-futures.md` and `key-flows/commodity-card-refresh.md`.
5. Do not infer canonical request fields, response fields, or fallback behavior from legacy code when a boundary file exists.

## Canonical Name Rule

Every type, event, symbol, query parameter, and response field defined in an integration-boundary file is canonical. Downstream agents must use those names exactly or include an explicit mapping table.
