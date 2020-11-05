package main

import (
	"fmt"
	"strings"
)

type opts struct {
	id     string
	config string // extra config provided to the analyzer
}

type args struct {
	exitcode int
	log      string
}

type result struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

type analyzer interface {
	run(args args) result
}

// ExitCodeAnalyzer is a ... TODO ...
type exitCodeAnalyzer struct {
	Category    string
	Subcategory string
}

func (a exitCodeAnalyzer) run(args args) result {
	if args.exitcode == 0 {
		return result{Category: "success", Subcategory: ""}
	}
	return result{Category: a.Category, Subcategory: a.Subcategory}
}

// SqApplyDiffsAnalyzer is a ... TODO ...
// type SqApplyDiffsAnalyzer struct{}

// func (a SqApplyDiffsAnalyzer) run(args args) result {
// 	// f(exitCode, log) -> result
// 	return result{Category: "baz", Subcategory: "moo"}
// }

func getAnalyzer(opts opts) (analyzer, error) {
	switch opts.id {
	case "exitcode":
		s := strings.Fields(opts.config)
		Category := s[0]
		Subcategory := s[1]
		return exitCodeAnalyzer{Category: Category, Subcategory: Subcategory}, nil
	// case "sq_apply_diffs":
	// 	return SqApplyDiffsAnalyzer{}, nil
	default:
		return nil, fmt.Errorf("invalid analyzer type %v", opts.id)
	}
}
