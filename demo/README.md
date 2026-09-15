# adr-rg demo

This small e-commerce example shows repository-governed Architecture Decision
Records for a database, API, and frontend.

## Quick start

Install the CLI from source with Go 1.26 or later:

~~~bash
go install github.com/averyfreeman/adr-repo-governance/cmd/adr@v0.2.0
~~~

The demo no longer carries a tar-download installer. From the parent checkout,
`make install` installs the CLI into the user-space layout; use
`make install-mac-layout` for the macOS application-data layout.

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
adr detect-bs --dir docs/adr --base origin/main
adr detect-bs --dir docs/adr --base origin/main --format json
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

## Try scaffolding

Generate a language-specific project in a separate directory without Git or
remote side effects:

~~~bash
adr scaffold --lang rust --dir /tmp/adr-rust-example
~~~

The generated `.adr-scaffold.yaml` contains editable version, commit, tag, and
push defaults.

## Layout

~~~text
demo/
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── docs/adr/
│   ├── README.md
│   ├── index.yaml
│   ├── templates/adr.md
│   └── *.md
└── src/
~~~

See the [main adr-rg documentation](../README.md) for the repository contract
and [demo ADR guidance](docs/adr/README.md) for the record format.
