# Demo governance instructions

This demo uses `adr-rg` to store and validate Architecture Decision Records.
Use the `adr` command as the repository's governance interface; the tool
identifies applicable decisions but does not prove semantic code compliance.

## Pre-flight

Install the CLI from source with Go 1.26 or later:

```bash
go install github.com/averyfreeman/adr-repo-governance/cmd/adr@v0.3.0
```

The demo intentionally has no tar-download installer scripts.

Before and after significant changes, validate the ADR set:

```bash
adr check --dir docs/adr --strict
```

When source files change, inspect applicable decisions with:

```bash
adr review --dir docs/adr --base origin/main
```

Read the matching ADRs before changing governed code. Treat their constraints
and invariants as repository requirements, and document a new decision when a
change cannot follow an existing one.

## Common commands

```bash
adr list --dir docs/adr
adr show --dir docs/adr ADR-0001
adr list --dir docs/adr --path "src/db/users.go"
adr list --dir docs/adr --tag database
adr new --dir docs/adr --tag tag1 --tag tag2 --path "affected/**" "Your Decision Title"
adr index --dir docs/adr
adr index --dir docs/adr --check
adr scaffold --dir /tmp/adr-rust-example rust
```

JSON is the integration format:

```bash
adr list --dir docs/adr --format json
```

Do not edit `index.yaml` manually; regenerate it with `adr index`. Follow the
required rationale pattern in the ADR README when adding or revising a record.
