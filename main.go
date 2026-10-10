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

	opts := cleaner.Options{DryRun: cfg.DryRun, Trash: cfg.Trash}

	past, future := "Deleted", "Would delete"
	if cfg.Trash {
		past, future = "Trashed", "Would trash"
	}
	verb := past
	if cfg.DryRun {
		verb = future
	}

	items, err := scanner.Scan(cfg.Root, cfg.Patterns, cfg.IgnoreDirs, cfg.MaxDepth)
	if err != nil {
		log.Printf("scan warnings:\n%v", err)
	}

	var dirs []string
	if cfg.EmptyDirs {
		dirs, err = scanner.EmptyDirs(cfg.Root, cfg.IgnoreDirs, cfg.MaxDepth, items)
		if err != nil {
			log.Printf("empty-dirs warnings:\n%v", err)
		}
	}

	if cfg.DryRun {
		for _, it := range items {
			fmt.Printf("[dry-run] %s (%d B)\n", it.Path, it.Size)
		}
	}

	files := cleaner.Remove(items, opts)
	report(cfg, past, files, "")

	failed := files.Failed
	fmt.Printf("%s: %d files, %d B freed, %d failed\n", verb, files.Deleted, files.Freed, files.Failed)

	if cfg.EmptyDirs {
		if cfg.DryRun {
			for _, d := range dirs {
				fmt.Printf("[dry-run] %s%c (empty dir)\n", d, os.PathSeparator)
			}
		}

		dres := cleaner.RemoveDirs(dirs, opts)
		report(cfg, past, dres, string(os.PathSeparator))

		failed += dres.Failed
		fmt.Printf("%s: %d empty dirs, %d failed\n", verb, dres.Deleted, dres.Failed)
	}

	if failed > 0 {
		os.Exit(1)
	}
}

func report(cfg config.Config, past string, res cleaner.Result, suffix string) {
	if cfg.Verbose && !cfg.DryRun {
		for _, p := range res.Removed {
			fmt.Printf("%s %s%s\n", past, p, suffix)
		}
	}

	for _, e := range res.Errors {
		log.Printf("error: %v", e)
	}
}
