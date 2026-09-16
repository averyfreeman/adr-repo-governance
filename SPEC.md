# adr-rg Specification

This document defines the file model and command contract for adr-rg, a
Git-native ADR repository-governance utility.

## Scope and non-goals

adr-rg creates, validates, indexes, queries, and scope-matches Architecture
Decision Records. It does not choose architecture, manage approval workflow,
implement source-code changes, or semantically prove that a change follows a
decision. Those responsibilities belong to humans, CI policy, or an architecture
assistant.

## Output formats

| Format | Use | Availability |
| --- | --- | --- |
| text | Human-readable command output | Default for CLI commands |
| json | Stable machine-readable command results | Commands that accept --format |
| yaml | Serialized index data | adr index only |

JSON is the canonical machine-readable integration format. Text is intended for
interactive use. The on-disk ADR index is YAML because it is generated metadata
that remains easy to review in Git.

Format names are case-insensitive. Command options must appear before a
positional title or identifier.

## ADR files

ADRs are stored in docs/adr/ unless --dir specifies another directory. Files use
the form NNNN-kebab-case-title.md and contain YAML frontmatter delimited by ---.

Required frontmatter:

~~~yaml
---
adr_id: ADR-NNNN
title: "Decision title"
status: proposed | adopted | rejected | deprecated | superseded
date: YYYY-MM-DD
---
~~~

Optional frontmatter includes:

~~~yaml
scope:
  paths:
    - "path/**"
tags: [database, api]
constraints:
  - "Rule that must be followed"
invariants:
  - "Property that must always hold"
supersedes: []
superseded_by: []
related_adrs: []
~~~

The body must contain these second-level headings:

1. Context
2. Decision
3. Alternatives Considered
4. Consequences

Strict validation also checks the repository’s configured rationale requirements.

## Generated index

adr index writes docs/adr/index.yaml by default. The file is derived from ADR
frontmatter and must not be edited manually. adr index --check compares the
existing entries with a freshly generated index while ignoring its timestamp.

## Commands

~~~text
adr [GLOBAL OPTIONS] <command> [OPTIONS]
adr init [OPTIONS]
adr new [OPTIONS] TITLE
adr list [OPTIONS]
adr show [OPTIONS] IDENTIFIER
adr check [OPTIONS]
adr index [OPTIONS]
adr review [--base REF] [OPTIONS]
adr scaffold [LANG] [OPTIONS]
adr version
~~~

Global `--dir`/`-d`, `--format`/`-o`, and `--json` options may appear before
the command. The same options remain available after each command for
compatibility. `adr detect-bs` is an alias for `adr review`.

### adr init

Creates the ADR directory and templates/adr.md when they do not exist, then
generates index.yaml from any ADR files already present. --dir changes the
target directory. JSON output contains adr_dir, template, and index paths.

### adr new

Creates the next numbered ADR with status proposed by default from the
initialized templates/adr.md when available, falling back to the embedded
template otherwise. It refreshes the index unless --no-index is provided.

| Option | Meaning |
| --- | --- |
| --dir PATH | ADR directory; default docs/adr |
| --tags TAGS, --tag TAG | Repeatable tags; comma-separated values remain supported |
| --paths GLOBS, --path GLOB | Repeatable scope globs; comma-separated values remain supported |
| --status STATUS | Initial lifecycle status |
| --no-index | Do not refresh the generated index |
| --format text\|json | Output format |

### adr scaffold

`scaffold` generates an offline project scaffold from the embedded language
catalog. It writes agent guidance, language-specific development guidance,
skills, `.gitignore`, README, and `.adr-scaffold.yaml` to the current directory
unless `--dir` is supplied.

~~~bash
adr scaffold rust
adr scaffold --lang go --dir ./new-project
~~~

Supported profiles include `typescript`, `javascript`, `go`, `rust`, `c`,
`cpp`, `html`, `haskell`, `lisp`/`clisp`, `lua`, `php`, `erlang`, `elixir`, and
`generic`, with common aliases such as `ts`, `js`, `rs`, `golang`, and `c++`.

The command never initializes Git, creates a remote, commits, tags, pushes, or
uses the network. `.adr-scaffold.yaml` records those defaults for later review
and editing. Existing files trigger an interactive collision prompt; `--force`
overwrites them without prompting.

| Option | Meaning |
| --- | --- |
| --lang LANG | Language profile or alias; may also be positional |
| --dir PATH | Target directory; default `.` |
| --force | Overwrite existing files without prompting |
| --format text\|json | Output format |

### adr list

Lists ADR metadata from the current ADR files. The generated index is not used
as a potentially stale cache; --tag matches an ADR containing any requested
tag, and --path matches a declared scope path.

| Option | Meaning |
| --- | --- |
| --dir PATH | ADR directory; default docs/adr |
| --status STATUS | Filter by lifecycle status |
| --tag TAG | Repeatable filter by tag; comma-separated values remain supported |
| --path PATH | Repeatable filter by scope-path match; comma-separated values remain supported |
| --format text\|json | Output format |

JSON output has the shape { "count": N, "adrs": [...] }. Each entry contains
adr_id, title, status, date, and file. Repeated tag or path filters are ORed
within their category; status and other categories are combined with them.

### adr show

Displays one ADR selected by ADR-NNNN, a numeric ID, or a filename.

| Option | Meaning |
| --- | --- |
| --dir PATH | ADR directory; default docs/adr |
| --format text\|json | Output format |

JSON output includes identity, status, date, tags, scope paths, constraints,
invariants, the extracted decision section, and filename.

### adr check

Validates every ADR in the selected directory. The former `adr check adr` form
remains a compatibility alias.

| Option | Meaning |
| --- | --- |
| --dir PATH | ADR directory; default docs/adr |
| --strict | Treat validation warnings as failures |
| --format text\|json | Output format |

JSON output contains valid, count, per-file results, and aggregate errors and
warnings.

### adr index

Generates the index unless --check is provided. In check mode it reports whether
the committed index matches current ADR metadata and never modifies the file.

| Option | Meaning |
| --- | --- |
| --dir PATH | ADR directory; default docs/adr |
| --check | Check freshness without writing |
| --format text\|json\|yaml | Output format |

### adr review

review runs git diff --name-only REF, loads ADRs from the selected
directory, and matches changed paths against each ADR’s scope.paths glob
patterns. If --base is omitted, the CLI uses an unambiguous `origin/HEAD`,
`origin/main`, or `main` reference when available. `detect-bs` remains a
compatibility alias.

~~~bash
adr review
adr review --base main --json
~~~

Only proposed and adopted ADRs are applicability candidates. Rejected,
deprecated, and superseded records remain available through adr list and
adr show as history but do not create active review signals.

The result reports:

- changed_files;
- applicable_adrs, including IDs, titles, matched paths, matched files,
  constraints, and invariants; and
- summary, including applicable ADR, constraint, and invariant counts plus
  aggregated constraints and invariants.

A successful match means that the decision should be reviewed. It does not mean
the implementation is compliant.

| Option | Meaning |
| --- | --- |
| --base REF | Optional Git ref passed to git diff --name-only; auto-detected when unambiguous |
| --dir PATH | ADR directory; default docs/adr |
| --format text\|json | Output format |

## Exit behavior

- 0: command completed successfully; non-strict validation warnings may exist;
- 1: command, Git, or input error; and
- 2: ADR validation failed, strict validation found warnings, or adr index
  --check found a stale index.

## Design principles

- Keep ADRs human-readable and version-controlled.
- Keep JSON stable for automation and AI-tool integration.
- Keep scope matching explicit and deterministic.
- Surface governance information without pretending to enforce semantics.
- Prefer small, testable components and minimal dependencies.
