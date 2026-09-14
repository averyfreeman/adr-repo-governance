package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunInitJSON(t *testing.T) {
	dir := t.TempDir()
	var output bytes.Buffer

	err := RunInit(&InitConfig{
		Dir: dir,
		Output: &Output{
			Format: FormatJSON,
			Writer: &output,
		},
	})
	if err != nil {
		t.Fatalf("RunInit() error = %v", err)
	}

	var result InitResult
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatalf("init output is not valid JSON: %v\n%s", err, output.String())
	}

	if result.ADRDir != dir {
		t.Errorf("ADRDir = %q, want %q", result.ADRDir, dir)
	}

	for name, path := range map[string]string{
		"template": result.Template,
		"index":    result.Index,
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s path %q is not present: %v", name, path, err)
		}
	}

	if want := filepath.Join(dir, "templates", "adr.md"); result.Template != want {
		t.Errorf("Template = %q, want %q", result.Template, want)
	}
}
