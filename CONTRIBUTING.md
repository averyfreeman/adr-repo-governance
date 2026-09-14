# Contributing to adr-rg

Contributions should keep `adr-rg` small, deterministic, and focused on ADR
repository governance. Architecture discovery, pattern selection, approval, and
implementation belong in the project’s architecture assistant rather than in this
CLI.

## Development

```bash
go test ./...
go vet ./...
go build ./cmd/adr
adr check adr --strict
adr index --check
```

Run `gofmt` on changed Go files. Prefer standard-library solutions and small,
testable functions. Return contextual errors rather than panicking.

## User-facing changes

Update the relevant documentation and add tests for new behavior. Changes to ADR
metadata, index generation, scope matching, or command output should include
fixtures covering both human-readable and JSON behavior where applicable.

Use conventional commit subjects such as:

```text
feat(cli): add ADR query filter
fix(index): detect stale scope metadata
docs(readme): clarify bs-detector behavior
```

Create or update an ADR when a change has long-term architectural consequences.
Run `adr check adr --strict` and `adr index --check` before submitting a change.
