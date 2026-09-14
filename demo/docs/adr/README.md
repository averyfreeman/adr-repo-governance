# Demo ADRs

These records document durable choices for the e-commerce example. Each one has
context, a decision, alternatives, consequences, constraints, invariants, and
scope paths.

## Current records

| ID | Title | Status | Scope |
| --- | --- | --- | --- |
| ADR-0001 | Use PostgreSQL for Persistence | adopted | `src/db/**`, migrations |
| ADR-0002 | REST API with Versioning | adopted | `src/api/**`, handlers |
| ADR-0003 | Redux for Frontend State Management | deprecated | frontend state |

## Commands

~~~bash
adr check adr --dir docs/adr --strict
adr list --dir docs/adr
adr list --dir docs/adr --path "src/db/users.go"
adr show --dir docs/adr ADR-0001
adr bs-detector --dir docs/adr --base main
adr index --dir docs/adr
adr index --dir docs/adr --check
~~~

Create a record with:

~~~bash
adr new --dir docs/adr --tags tag1 --paths "affected/**" "Your Decision"
~~~

Complete every required section before the strict check. Use `--format json`
for programmatic consumers.

## Rationale pattern

Adopted options must explain both sides of the choice:

~~~markdown
### [Option]: Adopted

**Adopted because:**
- Concrete reason tied to the decision drivers

**Adopted despite:**
- Known trade-off accepted by the project
~~~

Rejected alternatives use the corresponding `Rejected because` and `Rejected
despite` headings. The generated `index.yaml` is maintained by `adr index`; do
not edit it directly.
