package app

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckResult is a single doctor diagnostic.
type CheckResult struct {
	Name string
	OK   bool
	Msg  string
}

// Doctor runs environment diagnostics and returns one result per check.
func Doctor(cfg Config) []CheckResult {
	results := []CheckResult{
		{Name: "config resolved", OK: cfg.DBPath != "", Msg: cfg.DBPath},
	}
	if cfg.DBPath != "" {
		dir := filepath.Dir(cfg.DBPath)
		info, err := os.Stat(dir)
		ok := err == nil && info.IsDir()
		msg := dir
		if !ok {
			msg = fmt.Sprintf("%s (missing, will be created)", dir)
		}
		results = append(results, CheckResult{Name: "data directory", OK: true, Msg: msg})
	}
	if term := os.Getenv("TERM"); term != "" {
		results = append(results, CheckResult{Name: "terminal", OK: true, Msg: term})
	} else {
		results = append(results, CheckResult{Name: "terminal", OK: false, Msg: "$TERM is empty"})
	}
	results = append(results, CheckResult{Name: "version", OK: true, Msg: Version})
	return results
}
