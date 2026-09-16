package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNewUsesInitializedTemplate(t *testing.T) {
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("creating template directory: %v", err)
	}

	template := `---
adr_id: ADR-NNNN
title: "Your Decision Title Here"
status: proposed
date: YYYY-MM-DD
---

# ADR-NNNN: Your Decision Title Here

## Context

Custom initialized context.

## Decision

### [Chosen Option]: Adopted

**Adopted because:**
- Custom rationale.

**Adopted despite:**
- Custom trade-off.

## Alternatives Considered

### [Alternative A]: Rejected

**Rejected because:**
- Custom rejection rationale.

**Rejected despite:**
- Custom alternative strength.

## Consequences

**Positive:**
- Custom positive consequence.

**Negative:**
- Custom negative consequence.
`
	templatePath := filepath.Join(templateDir, "adr.md")
	if err := os.WriteFile(templatePath, []byte(template), 0644); err != nil {
		t.Fatalf("writing template: %v", err)
	}

	result, err := RunNew(&NewConfig{
		Title:  "Use initialized template",
		Dir:    dir,
		Tags:   []string{"go"},
		Paths:  []string{"internal/**"},
		Status: "proposed",
		Output: &Output{Format: FormatText, Writer: &bytes.Buffer{}},
	})
	if err != nil {
		t.Fatalf("RunNew() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, result.File))
	if err != nil {
		t.Fatalf("reading created ADR: %v", err)
	}
	text := string(content)
	for _, want := range []string{
		"Custom initialized context.",
		"ADR-0001: Use initialized template",
		"tags:",
		"go",
		"internal/**",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("created ADR does not contain %q:\n%s", want, text)
		}
	}
}
