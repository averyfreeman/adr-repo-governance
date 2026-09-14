# Getting started

adr-rg is a small Go CLI for creating, validating, indexing, and querying
Architecture Decision Records.

## Install

From a checkout:

~~~bash
./scripts/install-adr.sh
export PATH="$PWD/tools/adr:$PATH"
adr version
~~~

Or install from source:

~~~bash
go install github.com/averyfreeman/adr-repo-governance/cmd/adr@v0.1.0
~~~

Windows users can run `scripts/install-adr.ps1` from PowerShell.

## Start a repository

~~~bash
adr init
~~~

This creates `docs/adr/`, the ADR template, and the generated `index.yaml`.
Use `--dir PATH` when the repository keeps its records elsewhere.

## Create a decision

~~~bash
adr new --tags database --paths "src/db/**" "Use PostgreSQL for persistence"
~~~

Finish the generated `Context`, `Decision`, `Alternatives Considered`, and
`Consequences` sections. Add constraints and invariants when the decision needs
to guide later implementation.

## Validate and index

~~~bash
adr check adr --strict
adr index
adr index --check
~~~

The strict check validates every ADR. `adr index` regenerates derived metadata;
`adr index --check` fails when the committed index is stale.

## Review a change

~~~bash
adr bs-detector --base main
adr bs-detector --base main --format json
~~~

The Bullshit Detector compares changed Git paths with the scopes on proposed and
adopted ADRs. A match tells you what to read; it is not a semantic compliance
verdict.

## A few useful queries

~~~bash
adr list --status adopted
adr list --tag database
adr list --path "src/db/users.go"
adr show ADR-0001
~~~

For the complete flag and output contract, see the [reference](../reference/README.md).
