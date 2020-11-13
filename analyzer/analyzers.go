package analyzer

import (
	"fmt"
	"strings"
)

type CliArgs struct {
	Type   string
	Config string // extra config provided to the analyzer
}

type Args struct {
	Exitcode        int
	StdoutStr       string
	StderrStr       string
	StdoutStderrStr string
}

type Result struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

type Analyzer interface {
	Run(args Args) Result
}

func GetAnalyzer(cliArgs CliArgs) (Analyzer, error) {
	switch cliArgs.Type {
	case "exitcode":
		s := strings.Fields(cliArgs.Config)
		Category := s[0]
		Subcategory := s[1]
		return ExitCodeAnalyzer{Category: Category, Subcategory: Subcategory}, nil
	default:
		return nil, fmt.Errorf("invalid analyzer type %v", cliArgs.Type)
	}
}
