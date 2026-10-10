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

	if cfg.DryRun {
		for _, it := range items {
			fmt.Printf("[dry-run] %s (%d B)\n", it.Path, it.Size)
		}
	}

	files := cleaner.Remove(items, cfg.DryRun)
	report(cfg, files, "")

	failed := files.Failed
	verb := "Deleted"
	if cfg.DryRun {
		verb = "Would delete"
	}
	fmt.Printf("%s: %d files, %d B freed, %d failed\n", verb, files.Deleted, files.Freed, files.Failed)

	if cfg.EmptyDirs {
		var gone []scanner.Candidate
		if cfg.DryRun {
			gone = items
		}

		dirs, err := scanner.EmptyDirs(cfg.Root, cfg.IgnoreDirs, gone)
		if err != nil {
			log.Printf("empty-dirs warnings:\n%v", err)
		}

		if cfg.DryRun {
			for _, d := range dirs {
				fmt.Printf("[dry-run] %s%c (empty dir)\n", d, os.PathSeparator)
			}
		}

		dres := cleaner.RemoveDirs(dirs, cfg.DryRun)
		report(cfg, dres, string(os.PathSeparator))

		failed += dres.Failed
		fmt.Printf("%s: %d empty dirs, %d failed\n", verb, dres.Deleted, dres.Failed)
	}

	if failed > 0 {
		os.Exit(1)
	}
}

func report(cfg config.Config, res cleaner.Result, suffix string) {
	if cfg.Verbose && !cfg.DryRun {
		for _, p := range res.Removed {
			fmt.Printf("deleted %s%s\n", p, suffix)
		}
	}

	for _, e := range res.Errors {
		log.Printf("error: %v", e)
	}
}
