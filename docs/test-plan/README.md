# Test Plan

> Verification strategy derived from the architecture contract. Defines success criteria for every feature.

**Owned by:** `@sdet-engineer`

## Structure

| File | Contents |
|------|----------|
| `features/ura-refresh.md` | One file per feature — test scenarios, edge cases, integration points, validation commands |
| `features/gold-oil-card-restore.md` | One file per feature — test scenarios, edge cases, integration points, validation commands |
| `smoke-test.md` | Clean-room smoke test procedure (stack-agnostic bootstrap + run verification) |

## Conventions

- **One feature per file** in `features/`.
- Each feature file includes scenarios, integration points, expected outcomes, and exact validation commands.
- **Smoke test** defines the clean-room bootstrap procedure for this Go/Fyne app.

## How to use this index

1. Read this README first to see all tested features.
2. Drill into the specific feature file you need.
3. When adding test coverage for a new feature, create the file in `features/` and add it to this README.

## Feature Test Index

| Feature | File | Coverage | Notes |
|---------|------|----------|-------|
| URA refresh restore | `features/ura-refresh.md` | External PMI bootstrap/search flow, parser outcomes, app refresh path | Proof metric requires fresh-session non-empty rows from both standalone and app refresh paths |
| Gold and oil card restore | `features/gold-oil-card-restore.md` | Live Yahoo futures boundary classification, gold/oil mapper fallback behavior, startup and Refresh card replacement | Proof metric requires both cards to render current live Yahoo data on startup and Refresh |
