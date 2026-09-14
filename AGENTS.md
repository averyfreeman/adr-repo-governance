# adr-rg Agent Instructions

This repository uses `adr-rg` to manage Architecture Decision Records. ADRs are
reviewable governance documents containing metadata, constraints, invariants, and
scope paths.

## Before changing files

Find decisions whose scope overlaps the planned change:

```bash
adr bs-detector --base main
```

Read applicable ADRs and determine whether the change requires a new decision or
an explicit superseding ADR. `bs-detector` is a review signal; it does not prove
semantic compliance.

## Before finishing

```bash
adr check adr --strict
adr index --check
go test ./...
go vet ./...
```

Update the index with `adr index` when ADR metadata changes. Do not edit
`index.yaml` manually.

## ADR workflow

Create a decision with:

```bash
adr new --tags tag1,tag2 --paths "affected/**" "Decision title"
```

Complete `Context`, `Decision`, `Alternatives Considered`, and `Consequences`,
then add concrete constraints and invariants where they are useful. Supersede an
old decision with a new ADR and reciprocal `supersedes`/`superseded_by` links.

## Repository standards

- Keep dependencies minimal and use standard Go conventions.
- Run `gofmt` and return wrapped errors rather than panicking.
- Do not commit secrets.
- Do not make breaking CLI or file-format changes without documenting the change
  in an ADR.
