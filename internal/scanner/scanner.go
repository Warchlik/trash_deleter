package scanner

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type Candidate struct {
	Path string
	Size int64
}

func Scan(root string, patterns, ignoreDirs []string) ([]Candidate, error) {
	var (
		found []Candidate
		errs  []error
	)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errs = append(errs, err)
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if path != root && slices.Contains(ignoreDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !isTrash(d.Name(), patterns) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		found = append(found, Candidate{Path: path, Size: info.Size()})

		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	return found, errors.Join(errs...)
}

func isTrash(name string, patterns []string) bool {
	for _, p := range patterns {
		ok, err := filepath.Match(p, name)
		if err == nil && ok {
			return true
		}
	}

	return false
}

// EmptyDirs returns directories under root (root itself excluded) that are
// empty, or would be empty once the files in removed are gone. Directories in
// ignoreDirs are never entered or returned. Children come before their parents,
// so the result can be removed in order with os.Remove.
func EmptyDirs(root string, ignoreDirs []string, removed []Candidate) ([]string, error) {
	gone := make(map[string]struct{}, len(removed))
	for _, c := range removed {
		gone[c.Path] = struct{}{}
	}

	var (
		dirs []string
		errs []error
		walk func(dir string) bool
	)

	walk = func(dir string) bool {
		entries, err := os.ReadDir(dir)
		if err != nil {
			errs = append(errs, err)
			return false
		}

		empty := true
		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			if _, ok := gone[p]; ok {
				continue
			}
			if e.IsDir() && !slices.Contains(ignoreDirs, e.Name()) && walk(p) {
				dirs = append(dirs, p)
				continue
			}
			empty = false
		}

		return empty
	}

	walk(root)

	return dirs, errors.Join(errs...)
}
