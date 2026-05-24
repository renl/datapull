# Clean-Room Smoke Test

## Prerequisites
- Clean worktree with no stale build artifacts required for verification

## Steps

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Run tests**
   ```bash
   go test ./...
   ```

3. **Run build/package**
   ```bash
   go build ./...
   ```

4. **Start application**
   ```bash
   go run .
   ```
   Expected: app starts without build/runtime errors and reaches the existing UI shell.

5. **Run format check**
   ```bash
   gofmt -l .
   ```
   Expected: no output.

6. **Standalone URA fetch verification**
   ```bash
   go run ./test_ura.go
   ```
   Expected: fetch completes from a fresh session and reports non-empty current URA rows.

7. **Commodity restore live verification**
   ```bash
   go test ./... -run 'TestFetch(Gold|Oil)PriceLive|TestYahooFinanceClassifier' -count=1
   ```
   Expected: gold and oil restore paths pass against current live Yahoo responses, including at least one partial/meta-only classification check.
