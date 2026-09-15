---
adr_id: ADR-0007
title: "JSON as Canonical Structured Output"
status: adopted
date: 2026-09-13
scope:
  paths:
    - "cmd/**"
    - "internal/**"
    - "docs/**"
    - "demo/**"
    - "SPEC.md"
tags:
  - format
  - cli
  - output
constraints:
  - JSON must be the canonical machine-readable output for normal commands
  - YAML remains the generated format for ADR index files
  - Human-oriented commands default to text output
  - Structured output must be deterministic and preserve ADR semantics
invariants:
  - JSON consumers do not need an alternate parser for CLI integration
  - The ADR index is generated from ADR metadata
  - detect-bs identifies decisions for review and does not claim semantic compliance
supersedes:
  - ADR-0006
superseded_by: []
related_adrs:
  - ADR-0002
---

# ADR-0007: JSON as Canonical Structured Output

## Context

adr-rg serves both humans working in a repository and automation that needs to
query ADR metadata. A compact alternate serializer can reduce tokens for some
uniform data, but every additional representation increases parser, validation,
documentation, and interoperability costs. The CLI must integrate with common
shell tools, JSON Schema-oriented consumers, AI Software Architect, and generic
language libraries without requiring a project-specific decoder.

## Decision

Use JSON as the canonical machine-readable output for normal `adr` commands.
Keep text as the default for human-oriented output and YAML for the generated
ADR index. Do not ship a TOON encoder, decoder, or active TOON output path.

### JSON as Canonical Output: Adopted

**Adopted because:**
- JSON is supported by standard libraries and common automation tools
- A single canonical structured representation reduces conversion ambiguity
- JSON remains compatible with MCP-style JSON-RPC integrations and plugin APIs
- Deterministic JSON is straightforward to diff, validate, cache, and archive
- The format does not require consumers to install a project-specific parser

**Adopted despite:**
- JSON can use more tokens than compact tabular representations
- JSON is less pleasant to read than the CLI's text presentation for humans
- Existing experiments with alternate structured output are no longer available

## Alternatives Considered

### TOON as Machine Output: Rejected

**Rejected because:**
- It adds a second parser and a conversion boundary without a negotiated transport contract
- Generic CLI, MCP, JSON Schema, and plugin tooling already expects JSON
- Compact output is useful selectively, but is not sufficient reason to make it canonical
- The project cannot claim interoperability beyond consumers that explicitly support TOON

**Rejected despite:**
- Potential token savings for large, uniform catalogs
- A line-oriented representation that can be readable in focused contexts
- Its ability to represent the JSON data model when implementations agree on the rules

### YAML for All Structured Output: Rejected

**Rejected because:**
- YAML has more parser and scalar-normalization edge cases than JSON
- It is not the existing contract for programmatic command output
- Unnecessarily broad YAML output would blur the role of the generated index

**Rejected despite:**
- Familiarity to repository maintainers
- Human readability for configuration and index inspection

## Consequences

JSON is the stable integration contract for AI Software Architect and other
tools. Consumers can use standard parsers, while humans retain concise text
output. The ADR index continues to be YAML because it is a generated repository
artifact, not a general command-response protocol. LLM-facing callers may
project JSON into another prompt representation at their own boundary, but
that projection is not an `adr-rg` wire contract.

The CLI no longer carries the maintenance and ambiguity cost of an alternate
serializer. Consumers that depended on the historical TOON default must switch
to `--format json` or to text output as appropriate.
