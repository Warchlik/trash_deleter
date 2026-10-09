package config

import "flag"

type Config struct {
	Root       string
	Patterns   []string
	IngoreDirs []string
	Yes        bool
}

var defaultPatters []string = []string{
	"*:Zone.Identifier",
	".DS_Store",
	"Thumbs.db",
	"*.tmp",
	"*~",
}

var defaultIgnoreDirs []string = []string{
	".git",
	"node_modules",
}

func Parse() Config {
	root := flag.String(
		"path",
		".",
		"Default path to parse",
	)

	yes := flag.Bool(
		"yes",
		false,
		"Default dry-run",
	)

	flag.Parse()

	return Config{
		Root:       *root,
		Yes:        *yes,
		Patterns:   defaultPatters,
		IngoreDirs: defaultIgnoreDirs,
	}
}
