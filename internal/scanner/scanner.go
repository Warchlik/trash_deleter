package scanner

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Candidate struct {
	Path string
	Size int64
}

func Scan(root string, patterns, ignoreDirs []string, maxDepth int) ([]Candidate, error) {
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
			if path == root {
				return nil
			}
			if slices.Contains(ignoreDirs, d.Name()) {
				return filepath.SkipDir
			}
			if maxDepth > 0 && depthOf(root, path) >= maxDepth {
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

func depthOf(root, path string) int {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return 0
	}

	return strings.Count(rel, string(filepath.Separator)) + 1
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

func EmptyDirs(root string, ignoreDirs []string, maxDepth int, removed []Candidate) ([]string, error) {
	gone := make(map[string]struct{}, len(removed))
	for _, c := range removed {
		gone[c.Path] = struct{}{}
	}

	var (
		dirs []string
		errs []error
		walk func(dir string, depth int) bool
	)

	walk = func(dir string, depth int) bool {
		entries, err := os.ReadDir(dir)
		if err != nil {
			errs = append(errs, err)
			return false
		}

		canEnter := maxDepth == 0 || depth+1 <= maxDepth

		empty := true
		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			if _, ok := gone[p]; ok {
				continue
			}
			if canEnter && e.IsDir() && !slices.Contains(ignoreDirs, e.Name()) && walk(p, depth+1) {
				dirs = append(dirs, p)
				continue
			}
			empty = false
		}

		return empty
	}

	walk(root, 0)

	return dirs, errors.Join(errs...)
}
