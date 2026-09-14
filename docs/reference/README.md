# adr-rg Reference

## Command surface

| Command | Purpose |
| --- | --- |
| `adr init [OPTIONS]` | Initialize the ADR directory and index |
| `adr new [OPTIONS] TITLE` | Create a numbered ADR |
| `adr list [OPTIONS]` | List ADRs with filters |
| `adr show [OPTIONS] IDENTIFIER` | Show one ADR |
| `adr check adr [OPTIONS]` | Validate ADR files |
| `adr index [OPTIONS]` | Generate or check the index |
| `adr bs-detector --base REF [OPTIONS]` | Find decisions relevant to changed paths |
| `adr version` | Print version and build metadata |

Put options before the positional title or identifier.

## Options

| Option | Commands | Meaning |
| --- | --- | --- |
| `--dir PATH` | all commands except version | ADR directory; default `docs/adr` |
| `--format text\|json` | normal commands | Human or machine-readable output |
| `--format text\|json\|yaml` | index | Index output format |
| `--status STATUS` | new, list | Initial or filtered lifecycle status |
| `--tag TAG` | list | Filter by tag |
| `--path PATH` | list | Filter by scope-path match |
| `--paths GLOBS` | new | Comma-separated scope globs |
| `--tags TAGS` | new | Comma-separated tags |
| `--no-index` | new | Skip index refresh |
| `--strict` | check adr | Treat warnings as failures |
| `--check` | index | Check freshness without writing |
| `--base REF` | bs-detector | Required Git diff base reference |

## Output formats

Text is the default for interactive use. JSON is the canonical machine-readable
format. YAML is reserved for the generated ADR index.

## JSON result shapes

- `init`: `adr_dir`, `template`, and `index`;
- `new`: ADR ID, title, filename, path, and number;
- `list`: `count` and `adrs`;
- `show`: ADR metadata, decision, governance fields, and filename;
- `check adr`: `valid`, `count`, per-file results, errors, and warnings;
- `index`: generation/check status, index path, and ADR count; and
- `bs-detector`: `changed_files`, `applicable_adrs`, and `summary`.

## ADR frontmatter

~~~yaml
adr_id: ADR-NNNN
title: "Decision title"
status: proposed | adopted | rejected | deprecated | superseded
date: YYYY-MM-DD
scope:
  paths: ["src/**"]
tags: [architecture]
constraints: []
invariants: []
supersedes: []
superseded_by: []
related_adrs: []
~~~

The body requires `Context`, `Decision`, `Alternatives Considered`, and
`Consequences` headings. Strict validation also enforces the configured rationale
pattern.

## Index

`docs/adr/index.yaml` is generated from frontmatter. Run `adr index` after
metadata changes and `adr index --check` in CI. Do not edit it directly.

## Bullshit Detector

`adr bs-detector` obtains changed paths from `git diff --name-only REF` and
matches each proposed or adopted ADR’s `scope.paths` globs. It reports matching
ADR IDs, titles, paths, files, constraints, and invariants. Rejected, deprecated,
and superseded records remain queryable history but do not create active signals.

The command identifies decisions requiring review. It does not prove semantic
compliance.
