# adr-rg demo

This small e-commerce example shows repository-governed Architecture Decision
Records for a database, API, and frontend.

## Quick start

Install the pinned, checksum-verified binary:

~~~bash
# macOS/Linux
./scripts/install-adr.sh
export PATH="$PWD/tools/adr:$PATH"

# Windows PowerShell
.scriptsinstall-adr.ps1
$env:PATH = "$PWD	oolsadr;$env:PATH"
~~~

Validate and inspect the decisions:

~~~bash
adr check adr --dir docs/adr --strict
adr list --dir docs/adr
adr list --dir docs/adr --path "src/db/users.go"
adr show --dir docs/adr ADR-0001
~~~

The path query finds records whose declared scope covers a file. It is an
applicability signal, not a semantic compliance proof.

## Review a change

From a branch with the demo’s baseline available:

~~~bash
adr bs-detector --dir docs/adr --base main
adr bs-detector --dir docs/adr --base main --format json
~~~

The first form is for people; the second is for scripts and tools.

## Create and maintain decisions

~~~bash
adr new --dir docs/adr --tags performance --paths "src/**" "Add Caching"
adr check adr --dir docs/adr --strict
adr index --dir docs/adr
adr index --dir docs/adr --check
~~~

Each ADR needs `Context`, `Decision`, `Alternatives Considered`, and
`Consequences`. Adopted and rejected options also follow the demo’s rationale
pattern.

## Layout

~~~text
demo/
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── scripts/
│   ├── install-adr.sh
│   └── install-adr.ps1
├── tools/
│   ├── adr.version
│   └── adr/
├── docs/adr/
│   ├── README.md
│   ├── index.yaml
│   ├── templates/adr.md
│   └── *.md
└── src/
~~~

See the [main adr-rg documentation](../README.md) for the repository contract
and [demo ADR guidance](docs/adr/README.md) for the record format.
