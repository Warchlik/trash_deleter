package scanner

import (
	"errors"
	"io/fs"
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
