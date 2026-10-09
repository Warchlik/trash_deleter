package cleaner

import (
	scanner "trash_deleter/internal/scanner"
)

type Cleaner struct {
	Deleted int
	Failed  int
	Freed   int64
	Errors  []error
}

func Remove(items []scanner.Candidate, dryRun bool) Cleaner
