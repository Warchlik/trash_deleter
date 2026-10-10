package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Config struct {
	Root       string
	Patterns   []string
	IgnoreDirs []string
	DryRun     bool
	Verbose    bool
	EmptyDirs  bool
}

var defaultPatterns = []string{
	"*:Zone.Identifier",
	".DS_Store",
	"Thumbs.db",
	"*.tmp",
	"*~",
}

var defaultIgnoreDirs = []string{
	".git",
	"node_modules",
}

func Parse() Config {
	root := flag.String(
		"path",
		".",
		"Directory to scan",
	)

	extra := flag.String(
		"patterns",
		"",
		"Extra comma-separated file name patterns added to the defaults (glob, e.g. \"*.log,backup_*,notes.bak\")",
	)

	ignore := flag.String(
		"ignore",
		"",
		"Extra comma-separated directory names to skip, added to the defaults (e.g. \"vendor,dist\")",
	)

	dryRun := flag.Bool(
		"dry-run",
		false,
		"Only list files that would be deleted, do not delete anything",
	)

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Print every deleted file and directory")
	flag.BoolVar(&verbose, "v", false, "Shorthand for -verbose")

	emptyDirs := flag.Bool(
		"empty-dirs",
		false,
		"Also remove empty directories, including ones left empty after deleting files",
	)

	flag.Parse()

	extraPatterns := splitList(*extra)
	for _, p := range extraPatterns {
		if _, err := filepath.Match(p, ""); err != nil {
			fmt.Fprintf(os.Stderr, "invalid pattern %q: %v\n", p, err)
			os.Exit(2)
		}
	}

	return Config{
		Root:       *root,
		DryRun:     *dryRun,
		Verbose:    verbose,
		EmptyDirs:  *emptyDirs,
		Patterns:   slices.Concat(defaultPatterns, extraPatterns),
		IgnoreDirs: slices.Concat(defaultIgnoreDirs, splitList(*ignore)),
	}
}

func splitList(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ",") {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}

	return out
}
