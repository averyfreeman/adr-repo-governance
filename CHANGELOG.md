# Changelog

All notable changes to `adr-rg` are documented here.

## [0.1.0] - 2026-09-13

### Changed

- Rebranded the repository-governance utility as `adr-rg`.
- Renamed the executable and Go command to `adr`.
- Replaced the former diff-applicability command with `adr bs-detector`.
- Made JSON the only normal machine-readable output format and retained YAML for
  generated indexes.

### Removed

- The private alternate structured-output encoder and decoder.
- Claude-specific ADR workflow assets that overlap with architecture-assistant
  responsibilities.
- The redundant narrative applicability command.

### Retained

- ADR creation, validation, indexing, querying, scope matching, CI integration,
  installers, and the demonstration project.
