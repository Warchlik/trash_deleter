package main

import (
	"fmt"
	"log"
	"os"

	"trash_deleter/internal/cleaner"
	"trash_deleter/internal/config"
	"trash_deleter/internal/scanner"
)

func main() {
	cfg := config.Parse()

	info, err := os.Stat(cfg.Root)
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsDir() {
		log.Fatalf("%s is not a directory", cfg.Root)
	}

	items, err := scanner.Scan(cfg.Root, cfg.Patterns, cfg.IgnoreDirs)
	if err != nil {
		log.Printf("scan warnings:\n%v", err)
	}

	for _, it := range items {
		if cfg.DryRun {
			fmt.Printf("[dry-run] %s (%d B)\n", it.Path, it.Size)
		}
	}

	res := cleaner.Remove(items, cfg.DryRun)

	for _, e := range res.Errors {
		log.Printf("error: %v", e)
	}

	verb := "Deleted"
	if cfg.DryRun {
		verb = "Would delete"
	}
	fmt.Printf("%s: %d files, %d B freed, %d failed\n", verb, res.Deleted, res.Freed, res.Failed)

	if res.Failed > 0 {
		os.Exit(1)
	}
}
