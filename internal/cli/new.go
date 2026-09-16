package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/averyfreeman/adr-repo-governance/internal/adr"
	"github.com/averyfreeman/adr-repo-governance/internal/index"
	"github.com/averyfreeman/adr-repo-governance/internal/validate"
)

// NewConfig holds configuration for the new command.
type NewConfig struct {
	Title   string
	Dir     string
	Tags    []string
	Paths   []string
	Owners  []string
	Status  string
	NoIndex bool
	Format  OutputFormat
	Output  *Output
}

// NewResult holds the result of creating a new ADR.
type NewResult struct {
	ADRID    string `json:"adr_id"`
	Title    string `json:"title"`
	File     string `json:"file"`
	FilePath string `json:"file_path"`
	Number   int    `json:"number"`
}

// RunNew creates a new ADR.
func RunNew(cfg *NewConfig) (*NewResult, error) {
	// Validate inputs
	if err := validate.ValidateTitle(cfg.Title); err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}
	if err := validate.ValidateTags(cfg.Tags); err != nil {
		return nil, fmt.Errorf("invalid tags: %w", err)
	}
	if err := validate.ValidateScopePaths(cfg.Paths); err != nil {
		return nil, fmt.Errorf("invalid paths: %w", err)
	}

	// Parse and validate status
	status, err := adr.ParseStatus(cfg.Status)
	if err != nil {
		return nil, err
	}

	// Find next number
	number, err := adr.FindNextNumber(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("finding next ADR number: %w", err)
	}

	// Generate filename
	filename := adr.GenerateFilename(number, cfg.Title)
	filePath := filepath.Join(cfg.Dir, filename)
	adrID := fmt.Sprintf("ADR-%04d", number)

	// Create frontmatter
	fm := &adr.Frontmatter{
		ADRID:        adrID,
		Title:        cfg.Title,
		Status:       status,
		Date:         time.Now().Format("2006-01-02"),
		Scope:        adr.Scope{Paths: cfg.Paths},
		Tags:         cfg.Tags,
		Constraints:  []string{},
		Invariants:   []string{},
		Supersedes:   []string{},
		SupersededBy: []string{},
		RelatedADRs:  []string{},
	}

	// Generate content
	content, err := generateADRContent(cfg.Dir, filename, fm)
	if err != nil {
		return nil, fmt.Errorf("generating ADR content: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, fmt.Errorf("creating directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("writing ADR file: %w", err)
	}

	result := &NewResult{
		ADRID:    adrID,
		Title:    cfg.Title,
		File:     filename,
		FilePath: filePath,
		Number:   number,
	}

	// Update index unless disabled
	if !cfg.NoIndex {
		if err := index.WriteToDir(cfg.Dir); err != nil {
			// Non-fatal: warn but don't fail
			cfg.Output.Error("warning: could not update index: %v", err)
		}
	}

	// Output result
	if cfg.Output.IsStructuredFormat() {
		_ = cfg.Output.PrintStructured(result)
	} else {
		cfg.Output.Success("Created %s", filePath)
		cfg.Output.Println("  ADR ID: %s", adrID)
		cfg.Output.Println("  Title:  %s", cfg.Title)
		cfg.Output.Println("  Status: %s", status)
	}

	return result, nil
}

func generateADRContent(adrDir, filename string, fm *adr.Frontmatter) (string, error) {
	templatePath := filepath.Join(adrDir, "templates", "adr.md")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("reading ADR template: %w", err)
		}
		templateContent = []byte(defaultTemplate)
	}

	_, body, err := adr.ExtractFrontmatter(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("parsing ADR template: %w", err)
	}

	body = strings.ReplaceAll(body, "ADR-NNNN", fm.ADRID)
	body = strings.ReplaceAll(body, "Your Decision Title Here", fm.Title)
	body = strings.ReplaceAll(body, "YYYY-MM-DD", fm.Date)

	fmStr, err := adr.SerializeFrontmatter(fm)
	if err != nil {
		return "", fmt.Errorf("serializing frontmatter: %w", err)
	}

	content := fmStr + strings.TrimLeft(body, "\n")
	parsed, err := adr.ParseADR(content, filename, filepath.Join(adrDir, filename))
	if err != nil {
		return "", fmt.Errorf("validating rendered ADR: %w", err)
	}
	validation := adr.Validate(parsed)
	if !validation.IsValid() {
		return "", fmt.Errorf("rendered ADR is invalid: %s", validation.Errors[0].Message)
	}

	return content, nil
}
