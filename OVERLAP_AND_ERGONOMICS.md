# ADR Command Overlap and Ergonomics

## Executive summary

The `adr` command surface is small enough to understand, and most commands have
distinct responsibilities. The main problems are not unnecessary business
capabilities; they are duplicated setup logic, one overly deep command path,
and inconsistent relationships between source ADR files and the generated
index.

The highest-value improvements are:

1. Make `adr check` the canonical spelling while retaining `adr check adr` as a
   compatibility alias.
2. Reconcile `adr init` and `adr new`, which currently create overlapping
   repository artifacts from different templates.
3. Make `adr list` explicitly safe against stale `index.yaml` data.
4. Centralize common `--dir` and `--format` options, and offer ergonomic
   aliases without removing the existing forms.
5. Give `detect-bs` a user-facing alias such as `review` or `impact`, while
   retaining the existing command for scripts and existing documentation.

The commands should not be collapsed into one generic operation. Creating an
ADR, querying one, validating its document, checking generated metadata, and
finding decisions relevant to a Git diff are different workflows. The best
direction is a smaller-feeling interface built from compatible aliases and
shared defaults, not removal of the core capabilities.

This report is based on direct source inspection and runtime command help. The
repository is not currently indexed by the optional codebase-memory graph, so
the findings do not rely on graph-based call-relationship claims.

## Current command surface

The top-level dispatcher exposes eight operational commands plus metadata and
help (`cmd/adr/main.go:23-57`).

| Command | Current behavior | Assessment |
| --- | --- | --- |
| `adr init` | Creates the ADR directory, an editable template, and an empty index. | Useful bootstrap convenience, but not required before `new`; currently overlaps with it. |
| `adr new TITLE` | Allocates the next number, writes a complete ADR, and normally regenerates the index. | Necessary core operation. |
| `adr list` | Lists metadata, using `index.yaml` when present and scanning ADR files only when the index cannot be loaded. | Necessary query operation, with a freshness risk. |
| `adr show ID` | Resolves an ID, number, or filename and reads the source ADR, including its decision section. | Necessary detail operation; not a duplicate of `list`. |
| `adr check adr` | Validates every ADR’s frontmatter, filename, required sections, and rationale rules. | Necessary quality gate; the nested spelling is unnecessary. |
| `adr index` | Generates the derived YAML index; `--check` verifies freshness without writing. | Necessary because the index is a committed repository artifact and CI contract. |
| `adr detect-bs --base REF` | Computes changed Git paths and finds proposed/adopted ADRs whose scope globs match them. | Necessary review workflow, but the name is opaque and the base flag is mandatory. |
| `adr scaffold --lang LANG` | Writes language-specific project guidance, skills, ignore files, README, and scaffold configuration. | Useful adjacent workflow, but not an ADR-directory command. |
| `adr version` | Prints version, commit, and build metadata. | Standard operational command; no meaningful overlap. |

The command documentation in `SPEC.md:79-86` and
`docs/reference/README.md:7-15` accurately describes this surface, but it also
repeats the same option model in several places.

## Functional and logical overlap

### `init` and `new`: overlapping bootstrap paths

`adr init` creates `templates/adr.md` and an empty `index.yaml`
(`internal/cli/init.go:22-58`). `adr new` independently creates the target
directory, generates its own document body, and writes a fresh index
(`internal/cli/new.go:37-112`, `internal/cli/new.go:127-180`). It never reads the
template produced by `init`.

This creates three inconsistencies:

- Running `adr init` is not required for the main creation workflow.
- Editing `docs/adr/templates/adr.md` does not change what `adr new` generates.
- The two templates have different shapes: `init` writes a generic alternatives
  table, while `new` writes the rationale headings required by the validator.

There is also a correctness edge case. If a directory already contains ADR
files but has no index, `adr init` writes an empty index rather than generating
one from those files. Since `adr list` trusts a loadable index, the newly created
empty index can temporarily hide existing ADRs.

`init` is therefore not strictly necessary for creating or querying ADRs. It is
only necessary as a convenience for materializing an editable template and
initial empty repository structure.

**Recommendation:** keep `init` for compatibility, but make its role explicit
and make the paths converge. The preferred design is for `new` to use the
initialized template when one exists, while retaining an embedded fallback for
repositories that skip `init`. `init` should generate a real index when ADR
files already exist, and should only create an empty index for a genuinely empty
ADR directory. If an editable template is not intended to be supported, the
alternative is to deprecate `init` and document `adr new` as the sole bootstrap
path.

### `new` and `index`: intentional composition with an escape hatch

`new` automatically calls `index.WriteToDir` unless `--no-index` is supplied
(`internal/cli/new.go:106-112`). `index` remains independently useful because
users can edit ADR metadata by hand, repair a missing index, or validate the
generated artifact in CI (`internal/cli/index.go:27-84`).

These are not duplicates, but `--no-index` exposes an advanced batching mode in
the primary creation command. It is useful when several files are being created
or transformed, but it adds a branch that most users never need.

**Recommendation:** keep automatic index updates as the default. Retain
`--no-index` for scripts and bulk workflows, but describe it as an advanced
option or de-emphasize it in normal help. Do not remove `index`; it is the
explicit repair and CI-facing command.

### `list`, `show`, and `detect-bs`: shared query concepts, different jobs

The three commands all expose ADR metadata, but their inputs differ:

- `list` filters a catalog by status, tag, or a path matched against declared
  scope paths (`internal/cli/list.go:37-125`).
- `show` resolves one identifier and reads the source document, including the
  full decision section (`internal/cli/show.go:34-109`, `internal/cli/show.go:112-148`).
- `detect-bs` starts with a Git diff, then matches changed files against
  proposed/adopted ADR scope globs (`internal/cli/check.go:174-225`).

The commands are logically related but not redundant. `list --path` answers
“which decisions declare this path?”, while `detect-bs --base REF` answers
“which decisions should be reviewed for this Git diff?”

The implementation would benefit from a shared repository/query layer for ADR
loading and scope matching, but the user-facing commands should remain separate.
The shared layer should also establish one source-of-truth policy for the index
and source documents.

### `list` and `index`: stale generated data

`list` uses `index.yaml` whenever it can parse it and falls back to source ADR
files only when loading the index returns an error (`internal/cli/list.go:41-86`).
It does not call `index.Check` or compare the index with current frontmatter.

That means a syntactically valid but stale index can return old titles, statuses,
tags, scope paths, or file names. `adr index --check` can detect this condition,
but ordinary users invoking `adr list` are not warned. `show` and `detect-bs`
read source ADRs directly, so the three commands can disagree during the same
working-tree state.

This is the most important functional overlap problem because it makes the
commands appear to disagree about repository state.

**Recommendation:** make source ADR files authoritative for interactive queries.
Either have `list` scan source files directly, or verify index freshness before
using it and fall back to source files when stale. Keep the generated index for
fast external consumption and CI, but do not let its presence silently override
newer source metadata.

### `check adr` and `index --check`: distinct validation layers

`check adr` loads every ADR and validates document structure, metadata, and
rationale rules (`internal/cli/check.go:48-103`). `index --check` compares the
committed generated index with freshly generated entries while ignoring its
timestamp (`internal/cli/index.go:31-55`, `internal/index/index.go:116-166`).

These checks are complementary, not duplicates:

- `check` asks whether the source decision documents are valid.
- `index --check` asks whether derived repository metadata is current.

The overlap is in the word “check,” not in the underlying responsibility. The
current `check adr` spelling is unnecessarily verbose because `adr` is its only
supported subcommand (`cmd/adr/main.go:295-313`).

**Recommendation:** make `adr check` the canonical command, accept `adr check
adr` as a compatibility alias, and keep `adr index --check` as the focused
generated-artifact check. A future convenience mode could offer `adr check
--index` or `adr verify` for both checks, but it should not replace the focused
commands used by CI.

### `init` and `scaffold`: adjacent setup, different scope

`init` operates inside an ADR directory. `scaffold` renders a project-level set
of guidance files and `.adr-scaffold.yaml` without performing Git or network
actions (`internal/scaffold/scaffold.go:148-213`). It does not create the ADR
records themselves.

They are both setup commands, but combining them would make language selection,
collision handling, and ADR-directory initialization harder to understand.
Keep them separate. If a future workflow wants both, provide documentation or
an explicit higher-level workflow rather than merging their responsibilities.

## Ergonomics and flag analysis

### Repeated common flags

The dispatcher reconstructs `--dir` and `--format` independently for nearly
every command (`cmd/adr/main.go:86-119`, `cmd/adr/main.go:180-219`,
`cmd/adr/main.go:251-254`, `cmd/adr/main.go:315-319`,
`cmd/adr/main.go:349-353`, `cmd/adr/main.go:394-399`). This causes:

- repeated parsing and error-handling code;
- repeated documentation tables;
- duplicated short aliases (`-d`, `-o`) that users must remember;
- different directory semantics for ADR commands versus scaffold.

The defaults are already sensible: `docs/adr` for ADR commands and `.` for
scaffold. The problem is repetition, not that every command exposes an override.

**Recommendation:** introduce shared command options internally and support a
consistent persistent-option form such as:

```text
adr --dir docs/adr --format json list
```

Keep the existing post-command spelling during a compatibility period. A
convenient `--json` alias would be clearer than `--format json` for one-off
automation, while `--format` remains the extensible form. `--dir` should remain
an explicit override; implicit directory discovery can be considered separately
because silently selecting a parent repository can surprise scripts.

### `check adr` requires one unnecessary token

The help output currently exposes a subcommand tree even though `adr` is the
only valid child. This is the clearest low-risk simplification:

```text
# Preferred
adr check --strict

# Compatibility spelling
adr check adr --strict
```

The implementation can dispatch both forms to the same handler without changing
validation behavior or exit codes.

### `detect-bs` is descriptive to maintainers but opaque to users

The command implements scope applicability review, but its name is an internal
or editorial phrase. The required `--base REF` is appropriate for deterministic
CI, yet it makes the common local invocation verbose:

```text
adr detect-bs --base main
```

**Recommendation:** add a clearer canonical alias such as `adr review` or
`adr impact`, retain `detect-bs` as a compatibility alias, and document the
command as “find decisions applicable to changed paths.” For local use, a safe
base resolver could choose `origin/HEAD`, `origin/main`, or `main` only when one
candidate exists; CI should continue to pass `--base` explicitly. If no unique
base exists, fail with a helpful message rather than guessing.

### Filter syntax is more restrictive than the internal model

`ListConfig` supports a slice of tags and treats filters as an “any matching tag”
query (`internal/cli/list.go:112-136`), but the CLI exposes only one string-valued
`--tag` flag (`cmd/adr/main.go:213-219`, `cmd/adr/main.go:231-234`). `new` and
`--paths`/`--tags` use comma-separated strings, which is compact but less
shell-friendly when values contain commas or when scripts want repeatable
arguments.

**Recommendation:** accept repeatable flags in addition to the current comma
syntax:

```text
adr list --tag database --tag storage
adr new --path 'internal/**' --path 'cmd/**' 'Decision title'
```

Document whether repeated list filters are OR or AND; the current implementation
uses OR for tags and AND across filter categories. Preserve comma-separated input
for compatibility.

### Required versus optional flags

Most commands already require very little:

- `new` requires only a title; tags, paths, status, and output format have
  defaults.
- `show` requires only an identifier and already accepts a number or filename.
- `list`, `check`, and `index` work with defaults.
- `detect-bs` requires a base because the comparison point materially changes
  the result.
- `scaffold` requires a language because it cannot infer a profile reliably.

The required `--base` and `--lang` flags should not simply disappear. They can
gain safe positional or auto-detected conveniences, but the explicit forms are
valuable in automation. `--force` should remain explicit because it changes
collision behavior and can overwrite files.

## Recommended target interface

This is a compatibility-preserving target, not an immediate breaking-change
proposal:

| Current interface | Preferred interface | Compatibility treatment |
| --- | --- | --- |
| `adr check adr [OPTIONS]` | `adr check [OPTIONS]` | Keep `check adr` as an alias. |
| `adr detect-bs --base REF` | `adr review --base REF` | Keep `detect-bs` as an alias; retain explicit `--base` in CI. |
| `adr scaffold --lang rust` | `adr scaffold rust` or current form | Accept positional language as shorthand; retain `--lang`. |
| `adr show ADR-0001` | `adr show 0001` | Already supported; use the shorter form in examples. |
| `adr --format json list` | Shared global output option | Continue accepting `list --format json`; optionally add `--json`. |
| `adr --dir PATH list` | Shared repository option | Continue accepting `list --dir PATH`. |
| `adr list --tag A --tag B` | Repeatable filters | Preserve existing single-value and comma-separated forms. |

The command names should remain stable until aliases and documentation have been
available for at least one release. Scripts should be able to detect the newer
spelling through help output or continue using the existing commands unchanged.

## Prioritized follow-up work

### Priority 1: correctness and low-risk ergonomics

1. Add `adr check` and retain `adr check adr`.
2. Make `list` source-consistent by detecting stale indexes or scanning source
   ADRs for interactive queries.
3. Fix `init` so an existing ADR set cannot be hidden by a newly created empty
   index.
4. Add command-level tests for the common lifecycle:
   `init` → `new` → `list` → `show` → `check` → `index --check`.

### Priority 2: shared interface improvements

1. Centralize common option parsing and output construction.
2. Add repeatable tag/path flags and document their combination semantics.
3. Add `--json` and shorter examples for numeric `show` identifiers.
4. Add the clearer review-oriented alias for `detect-bs`.

### Priority 3: template and configuration cleanup

1. Decide whether the initialized template is a supported customization point.
2. If yes, make `new` load it with a validated embedded fallback.
3. If no, deprecate the materialized template and simplify `init` to a narrow
   repository bootstrap or remove it in a future major release.
4. Consider repository-level configuration or safe directory discovery only
   after the explicit `--dir` contract is covered by compatibility tests.

## Test and acceptance scenarios

Any future implementation should cover these cases:

- `adr check --strict` and `adr check adr --strict` produce equivalent results
  and exit codes.
- `adr init` in an empty directory creates a valid empty index; `adr init` in a
  directory containing ADR files generates an index containing those files.
- Editing an ADR after index generation does not make `adr list` silently return
  stale metadata.
- `adr new` uses the selected template policy and continues to update the index
  by default.
- `adr index --check` remains a read-only freshness check with CI-friendly exit
  status.
- `adr review --base REF` and `adr detect-bs --base REF` return identical JSON
  shapes and applicability results.
- Global and post-command forms of `--dir` and `--format` resolve to the same
  configuration.
- Repeatable and comma-separated tag/path filters behave consistently.
- `--force` remains required for non-interactive scaffold overwrites.
- Help output clearly distinguishes source validation, index freshness, and
  changed-path applicability review.

## Conclusion

The command set is broadly necessary, but the current interface makes related
operations look more fragmented than they are. `init` is the only command whose
current behavior is substantially redundant with another command, and `list`
has the most important consistency risk because it can trust stale derived data.

The strongest near-term result is a compatibility layer that makes the common
paths read naturally:

```text
adr new "Decision title"
adr list --status adopted
adr show 0007
adr check --strict
adr index --check
adr review --base main
```

This keeps the repository’s explicit governance model while removing avoidable
verbosity and making the distinction between document validity, generated-index
freshness, and changed-path review visible to users.
