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
	"vendor",
	"venv", ".venv",
}

func Parse() Config {
	root := flag.String(
		"path",
		".",
		"Directory to scan",
	)

	extraPatterns := flag.String(
		"patterns",
		"",
		"Extra comma-separated file name patterns added to the defaults (glob, e.g. \"*.log,backup_*,notes.bak\")",
	)

	extraIgnore := flag.String(
		"ignore",
		"",
		"Extra comma-separated directory names to skip, added to the defaults (e.g. \"vendor,dist\")",
	)

	dryRun := flag.Bool(
		"dry-run",
		false,
		"Only list files that would be deleted, do not delete anything",
	)

	flag.Parse()

	patternsArr := splitList(*extraPatterns)
	for _, p := range patternsArr {
		if _, err := filepath.Match(p, ""); err != nil {
			fmt.Fprintf(os.Stderr, "invalid pattern %q: %v\n", p, err)
			os.Exit(2)
		}
	}

	return Config{
		Root:       *root,
		DryRun:     *dryRun,
		Patterns:   slices.Concat(defaultPatterns, patternsArr),
		IgnoreDirs: slices.Concat(defaultIgnoreDirs, splitList(*extraIgnore)),
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
