package cleaner

import (
	"os"

	"trash_deleter/internal/scanner"
)

type Result struct {
	Deleted int
	Failed  int
	Freed   int64
	Errors  []error
}

func Remove(items []scanner.Candidate, dryRun bool) Result {
	var res Result

	for _, it := range items {
		if !dryRun {
			if err := os.Remove(it.Path); err != nil {
				res.Failed++
				res.Errors = append(res.Errors, err)
				continue
			}
		}

		res.Deleted++
		res.Freed += it.Size
	}

	return res
}
