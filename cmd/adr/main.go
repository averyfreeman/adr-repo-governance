// Package main provides the entry point for the adr CLI.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/averyfreeman/adr-repo-governance/internal/cli"
	"github.com/averyfreeman/adr-repo-governance/internal/scaffold"
)

// Version information, set via ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const defaultADRDir = "docs/adr"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		runInit(os.Args[2:])
	case "new":
		runNew(os.Args[2:])
	case "index":
		runIndex(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	case "show":
		runShow(os.Args[2:])
	case "check":
		runCheck(os.Args[2:])
	case "detect-bs":
		runDetectBS(os.Args[2:])
	case "scaffold":
		runScaffold(os.Args[2:])
	case "version":
		printVersion()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`adr-rg - Git-native ADR governance

Usage:
  adr <command> [options]

Commands:
  init          Initialize ADR directory structure
  new           Create a new ADR
  index         Generate/update the ADR index
  list          List ADRs with optional filters
  show          Display details of an ADR
  check         Validate ADRs
  detect-bs     Detect ADRs relevant to changed files
  scaffold      Generate a language-specific project scaffold
  version       Show version information
  help          Show this help message

Run 'adr <command> -h' for more information on a command.`)
}

func printVersion() {
	fmt.Printf("adr version %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built:  %s\n", date)
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.InitConfig{
		Dir:    *dir,
		Output: cli.NewOutput(outputFormat),
	}

	if err := cli.RunInit(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	tags := stringFlag(fs, "tags", "t", "", "Comma-separated tags")
	paths := stringFlag(fs, "paths", "p", "", "Comma-separated scope paths (globs)")
	status := stringFlag(fs, "status", "s", "proposed", "Initial status")
	noIndex := boolFlag(fs, "no-index", "n", false, "Skip updating index")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: adr new [options] <title>")
		fmt.Println()
		fmt.Println("Create a new ADR with the given title.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: title is required")
		fs.Usage()
		os.Exit(1)
	}

	title := strings.Join(fs.Args(), " ")

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var tagList []string
	if *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			tagList = append(tagList, strings.TrimSpace(t))
		}
	}

	var pathList []string
	if *paths != "" {
		for _, p := range strings.Split(*paths, ",") {
			pathList = append(pathList, strings.TrimSpace(p))
		}
	}

	cfg := &cli.NewConfig{
		Title:   title,
		Dir:     *dir,
		Tags:    tagList,
		Paths:   pathList,
		Status:  *status,
		NoIndex: *noIndex,
		Format:  outputFormat,
		Output:  cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunNew(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runIndex(args []string) {
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	check := boolFlag(fs, "check", "c", false, "Check if index is up-to-date (don't modify)")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json|yaml)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseIndexOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.IndexConfig{
		Dir:    *dir,
		Check:  *check,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	_, err = cli.RunIndex(cfg)
	if err != nil {
		if *check {
			os.Exit(2) // Lint failure exit code
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	status := stringFlag(fs, "status", "s", "", "Filter by status")
	tag := stringFlag(fs, "tag", "t", "", "Filter by tag")
	path := stringFlag(fs, "path", "p", "", "Filter by scope path match")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var tags []string
	if *tag != "" {
		tags = []string{*tag}
	}

	cfg := &cli.ListConfig{
		Dir:    *dir,
		Status: *status,
		Tags:   tags,
		Path:   *path,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunList(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runShow(args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: adr show [options] <ADR-ID|number|filename>")
		fmt.Println()
		fmt.Println("Display details of a specific ADR.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: ADR identifier is required")
		fs.Usage()
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.ShowConfig{
		ID:     fs.Arg(0),
		Dir:    *dir,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunShow(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runCheck(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: adr check adr [options]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Subcommands:")
		fmt.Fprintln(os.Stderr, "  adr       Validate ADR files")
		os.Exit(1)
	}

	subCmd := args[0]

	switch subCmd {
	case "adr":
		runCheckADR(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown check subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

func runCheckADR(args []string) {
	fs := flag.NewFlagSet("check adr", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	strict := boolFlag(fs, "strict", "s", false, "Treat warnings as errors (fail on missing rationale pattern)")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.CheckADRConfig{
		Dir:    *dir,
		Strict: *strict,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	result, err := cli.RunCheckADR(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !result.Valid {
		os.Exit(2) // Lint failure exit code
	}
}

func runDetectBS(args []string) {
	fs := flag.NewFlagSet("detect-bs", flag.ExitOnError)
	dir := stringFlag(fs, "dir", "d", defaultADRDir, "ADR directory path")
	base := stringFlag(fs, "base", "b", "", "Base ref for git diff (required)")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: adr detect-bs --base <ref> [options]")
		fmt.Println()
		fmt.Println("Detect ADRs whose documented scope overlaps files changed since <base>.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *base == "" {
		fmt.Fprintln(os.Stderr, "error: --base is required")
		fs.Usage()
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.DetectBSConfig{
		Dir:    *dir,
		Base:   *base,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunDetectBS(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runScaffold(args []string) {
	fs := flag.NewFlagSet("scaffold", flag.ExitOnError)
	language := stringFlag(fs, "lang", "l", "", "Language template (required)")
	dir := stringFlag(fs, "dir", "d", ".", "Target directory")
	force := boolFlag(fs, "force", "F", false, "Overwrite existing files without prompting")
	format := stringFlag(fs, "format", "o", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: adr scaffold --lang <language> [options]")
		fmt.Println()
		fmt.Println("Generate an offline, language-specific project scaffold.")
		fmt.Println()
		fmt.Println("Available languages:", strings.Join(scaffold.Languages(), ", "))
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	result, err := scaffold.Run(&scaffold.Config{
		Language: *language,
		Dir:      *dir,
		Force:    *force,
		Input:    os.Stdin,
		Output:   os.Stderr,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	output := cli.NewOutput(outputFormat)
	if output.IsStructuredFormat() {
		_ = output.PrintStructured(result)
		return
	}
	output.Success("Scaffolded %s project in %s", result.Language, result.Directory)
	for _, file := range result.Files {
		output.Println("  Created %s", file)
	}
}

func stringFlag(fs *flag.FlagSet, name, short, value, usage string) *string {
	result := fs.String(name, value, usage)
	fs.StringVar(result, short, value, usage+" (short form)")
	return result
}

func boolFlag(fs *flag.FlagSet, name, short string, value bool, usage string) *bool {
	result := fs.Bool(name, value, usage)
	fs.BoolVar(result, short, value, usage+" (short form)")
	return result
}
