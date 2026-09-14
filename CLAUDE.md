# adr-rg Repository Guidance

Use `adr-rg` as the repository-governance utility. Architecture reasoning and
approval belong to the project’s architecture assistant; this repository provides
the durable ADR store and deterministic checks.

Before changing code:

```bash
adr bs-detector --base main --format json
```

Before finishing:

```bash
adr check adr --strict
adr index --check
go test ./...
```

Use `adr index` after changing ADR metadata. Treat `bs-detector` results as
decisions requiring review, not as semantic compliance proof.
