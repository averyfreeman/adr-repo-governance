# Demo governance instructions

This demo uses `adr-rg` to store and validate Architecture Decision Records.
Use the `adr` command as the repository's governance interface; the tool
identifies applicable decisions but does not prove semantic code compliance.

## Pre-flight

Install the pinned demo binary when needed:

```bash
# macOS/Linux
./scripts/install-adr.sh
export PATH="$PWD/tools/adr:$PATH"

# Windows PowerShell
.\scripts\install-adr.ps1
$env:PATH = "$PWD\tools\adr;$env:PATH"
```

Before and after significant changes, validate the ADR set:

```bash
adr check adr --dir docs/adr --strict
```

When source files change, inspect applicable decisions with:

```bash
adr bs-detector --dir docs/adr --base origin/main
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
adr new --dir docs/adr --tags tag1,tag2 --paths "affected/**" "Your Decision Title"
adr index --dir docs/adr
adr index --dir docs/adr --check
```

JSON is the integration format:

```bash
adr list --dir docs/adr --format json
```

Do not edit `index.yaml` manually; regenerate it with `adr index`. Follow the
required rationale pattern in the ADR README when adding or revising a record.
