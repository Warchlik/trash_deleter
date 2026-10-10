package cleaner

import (
	"fmt"
	"os"

	"github.com/Bios-Marcel/wastebasket/v2"

	"trash_deleter/internal/scanner"
)

type Options struct {
	DryRun bool
	Trash  bool
}

type Result struct {
	Deleted int
	Failed  int
	Freed   int64
	Removed []string
	Errors  []error
}

func Remove(items []scanner.Candidate, opts Options) Result {
	var res Result

	for _, it := range items {
		if !opts.DryRun {
			if err := removeFile(it.Path, opts.Trash); err != nil {
				res.Failed++
				res.Errors = append(res.Errors, err)
				continue
			}
		}

		res.Deleted++
		res.Freed += it.Size
		res.Removed = append(res.Removed, it.Path)
	}

	return res
}

func RemoveDirs(dirs []string, opts Options) Result {
	var res Result

	for _, d := range dirs {
		if !opts.DryRun {
			if err := removeDir(d, opts.Trash); err != nil {
				res.Failed++
				res.Errors = append(res.Errors, err)
				continue
			}
		}

		res.Deleted++
		res.Removed = append(res.Removed, d)
	}

	return res
}

func removeFile(path string, trash bool) error {
	if trash {
		return wastebasket.Trash(path)
	}

	return os.Remove(path)
}

func removeDir(path string, trash bool) error {
	if !trash {
		return os.Remove(path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Name() != ".DS_Store" {
			return fmt.Errorf("%s: directory not empty", path)
		}
	}

	return wastebasket.Trash(path)
}
