# adr-rg ADRs

This directory contains the architectural decisions for `adr-rg` itself.

## Current decisions

| ID | Decision | Status |
| --- | --- | --- |
| ADR-0001 | Adopt Go for the adr-rg CLI | adopted |
| ADR-0002 | Markdown ADRs with YAML frontmatter | adopted |
| ADR-0003 | Repository layout and generated index | adopted |
| ADR-0004 | GoReleaser and GitHub Actions releases | adopted |
| ADR-0005 | Explicit rationale for durable decisions | adopted |
| ADR-0006 | Historical structured-output experiment | superseded |
| ADR-0007 | JSON as canonical structured output | adopted |

## Working with ADRs

```bash
adr new --tags tooling --paths "internal/**" "Decision title"
adr check adr --strict
adr index
adr bs-detector --base main
```

`index.yaml` is generated; use `adr index` rather than editing it manually. ADR
history is preserved when a decision is superseded.

## Required sections

Every ADR contains `Context`, `Decision`, `Alternatives Considered`, and
`Consequences`. Frontmatter may additionally declare constraints, invariants, and
scope paths used by `bs-detector`.
