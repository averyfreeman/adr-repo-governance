# adr-rg Repository Guidance

Use `adr-rg` as the repository-governance utility. Architecture reasoning and
approval belong to the project’s architecture assistant; this repository provides
the durable ADR store and deterministic checks.

Before changing code:

```bash
adr review --base main --format json
```

Before finishing:

```bash
adr check --strict
adr index --check
make test
make vet
make build
```

Use `adr index` after changing ADR metadata. Treat `adr review` results as
decisions requiring review, not as semantic compliance proof. `detect-bs`
remains a compatibility alias.
