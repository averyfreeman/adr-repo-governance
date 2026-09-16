# Demo repository governance

Follow [AGENTS.md](AGENTS.md) for the complete, tool-neutral workflow. The
commands below are the normal `adr-rg` integration points for this demo:

```bash
adr check --dir docs/adr --strict
adr review --dir docs/adr --base origin/main
adr list --dir docs/adr --path "src/api/**"
adr show --dir docs/adr ADR-0001
adr index --dir docs/adr --check
```

Use `adr new` for new decisions, preserve required metadata and rationale
sections, and use `--format json` when another tool consumes the result. Do
not edit the generated index directly.
