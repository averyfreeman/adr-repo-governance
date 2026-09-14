# Tool integration

adr-rg gives coding tools and plugins a small, stable JSON surface for reading
decisions before a change and reviewing them afterward.

## Before implementation

~~~bash
adr list --format json
adr show --format json ADR-0001
adr bs-detector --base main --format json
~~~

Use `list` for discovery, `show` for the full decision context, and
`bs-detector` for path-scoped review context. Detector results include:

- `changed_files`;
- `applicable_adrs`;
- matched scope paths and files; and
- the applicable constraints and invariants.

The caller still has to read the ADR and evaluate the implementation. A path
match is a prompt to review, not proof of a violation.

## After implementation

~~~bash
adr bs-detector --base main --format json
adr check adr --strict
adr index --check
~~~

Prefer JSON for scripts and plugin calls, text for interactive review, and
version-controlled ADRs plus the generated index for durable context.
