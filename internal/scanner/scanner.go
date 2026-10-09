package scanner

import (
	"io/fs"
	"path/filepath"
	"slices"
)

type Candidate struct {
	Path string
	Size int64
}

func Scan(root string, patterns, ignoreDirs []string) ([]Candidate, error) {
	var found []Candidate

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root && slices.Contains(ignoreDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !isTrush(d.Name(), patterns) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		found = append(found, Candidate{Path: path, Size: info.Size()})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return found, nil
}

func isTrush(name string, patterns []string) bool {
	for _, p := range patterns {
		ok, err := filepath.Match(p, name)
		if err == nil && ok {
			return true
		}
	}

	return false
}
