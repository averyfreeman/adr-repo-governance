# Changelog

All notable changes to `adr-rg` are documented here.

## [0.2.0] - 2026-09-15

### Added

- Added `adr scaffold --lang LANG` with embedded language templates.
- Added Makefile targets for formatting, testing, vetting, building, and installing.
- Added editable `.adr-scaffold.yaml` Git workflow defaults and collision prompts.

### Changed

- Raised the supported Go toolchain to 1.26 or later.
- Made `adr detect-bs` the canonical scope-review command.

### Removed

- Removed broken tar-download installer scripts and obsolete version-pin files.

## [0.1.0] - 2026-09-13

### Changed

- Rebranded the repository-governance utility as `adr-rg`.
- Renamed the executable and Go command to `adr`.
- Replaced the former diff-applicability command with `adr detect-bs`.
- Made JSON the only normal machine-readable output format and retained YAML for
  generated indexes.

### Removed

- The private alternate structured-output encoder and decoder.
- Claude-specific ADR workflow assets that overlap with architecture-assistant
  responsibilities.
- The redundant narrative applicability command.

### Retained

- ADR creation, validation, indexing, querying, scope matching, CI integration,
  and the demonstration project.
