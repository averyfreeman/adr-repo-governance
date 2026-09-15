---
adr_id: ADR-NNNN
title: "Your Decision Title Here"
status: proposed
date: YYYY-MM-DD
scope:
  paths:
    - "path/to/affected/**"
tags:
  - tag1
  - tag2
constraints:
  - "First constraint that MUST be followed"
  - "Second constraint"
invariants:
  - "Property that must always hold true"
supersedes: []
superseded_by: []
related_adrs: []
---

# ADR-NNNN: Your Decision Title Here

## Context

Describe the context and background that led to this decision:
 1. What problem are we solving?
 2. Why is solving this problem important?
 3. Is it part of a sequenced process (aka a dependency)?
 4. What are the most pertinent internal attributes?
 5. What are the most influential externalities?
 6. Consider the perceived future benefit of solving this problem vs. the null hypothesis (aka doing nothing). Are you sure you want to continue?

Decision drivers:
- Key driver 1 that influenced the decision
- Key driver 2
- Key driver 3

## Decision

State the decision clearly and concisely.

### [Chosen Option]: Adopted

**Adopted because:**
- Clear, concrete reason why this option was chosen
- Tie reasons to decision drivers above
- Technical, operational, or strategic justification

**Adopted despite:**
- Known downside or trade-off we consciously accepted
- Cost or weakness compared to alternatives
- Risk we are taking on

## Alternatives Considered

### [Alternative A]: Rejected

**Rejected because:**
- Clear, concrete reason why this option was not chosen
- Technical, organizational, or strategic reason
- How it failed to meet decision drivers

**Rejected despite:**
- Legitimate strength of this option
- Benefit that made it attractive
- Reason it was seriously considered

### [Alternative B]: Rejected

**Rejected because:**
- Reason 1
- Reason 2

**Rejected despite:**
- Strength 1
- Strength 2

## Consequences

**Positive:**
- First positive consequence
- Second positive consequence

**Negative:**
- First negative consequence (and mitigation if any)
- Second negative consequence

## Agent Guidance

_Optional section. Include specific instructions for AI agents working in the affected scope._

When working in this area:
- Specific guidance point 1
- Specific guidance point 2
