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
| `adr detect-bs --base REF [OPTIONS]` | Find decisions relevant to changed paths |
| `adr scaffold --lang LANG [OPTIONS]` | Generate a language-specific project scaffold |
| `adr version` | Print version and build metadata |

Put options before the positional title or identifier.

## Options

| Option | Commands | Meaning |
| --- | --- | --- |
| `--dir PATH` | ADR commands | ADR directory; default `docs/adr` |
| `--dir PATH` | scaffold | Target project directory; default `.` |
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
| `--base REF` | detect-bs | Required Git diff base reference |
| `--lang LANG` | scaffold | Required language template name or alias |
| `--force` | scaffold | Overwrite existing scaffold files without prompting |

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
- `detect-bs`: `changed_files`, `applicable_adrs`, and `summary`; and
- `scaffold`: selected language, target directory, generated files, and config path.

Short aliases are available for common options: `-d` for `--dir`, `-o` for
`--format`, `-b` for `--base`, `-l` for `--lang`, and `-F` for `--force`.

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

`adr detect-bs` obtains changed paths from `git diff --name-only REF` and
matches each proposed or adopted ADR’s `scope.paths` globs. It reports matching
ADR IDs, titles, paths, files, constraints, and invariants. Rejected, deprecated,
and superseded records remain queryable history but do not create active signals.

The command identifies decisions requiring review. It does not prove semantic
compliance.

## Scaffold

`adr scaffold --lang LANG` writes the default agent files, language guidance,
`.gitignore`, README, skills, and `.adr-scaffold.yaml` into the target directory.
Supported profiles include TypeScript, JavaScript, Go, Rust, C, C++, HTML,
Haskell, Common Lisp/CLISP, Lua, PHP, Erlang, Elixir, and generic projects.

Scaffolding performs no Git or network actions. Existing files are handled by an
interactive overwrite, skip, rename, or alternate-path prompt. Use `--force` to
overwrite all collisions without prompting.
