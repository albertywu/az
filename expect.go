package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// CliArgs is a ... TODO ...
type CliArgs struct {
	Type      string
	Args      string // extra args provided to the analyzer
	OutputDir string // directory to save failure analysis
}

// AnalyzerArgs is a ... TODO ...
type AnalyzerArgs struct {
	Exitcode int
	Log      string
}

// AnalyzerResult is a ... TODO ...
type AnalyzerResult struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

// Analyzer is a thing with an analyze function
type Analyzer interface {
	analyze(args AnalyzerArgs) AnalyzerResult
}

// ExitCodeAnalyzer is a ... TODO ...
type ExitCodeAnalyzer struct {
	Category    string
	Subcategory string
}

func (a ExitCodeAnalyzer) analyze(args AnalyzerArgs) AnalyzerResult {
	if args.Exitcode == 0 {
		return AnalyzerResult{Category: "success", Subcategory: ""}
	}
	return AnalyzerResult{Category: a.Category, Subcategory: a.Subcategory}
}

// SqApplyDiffsAnalyzer is a ... TODO ...
type SqApplyDiffsAnalyzer struct{}

func (a SqApplyDiffsAnalyzer) analyze(args AnalyzerArgs) AnalyzerResult {
	// f(exitCode, log) -> analyzerResult
	return AnalyzerResult{Category: "baz", Subcategory: "moo"}
}

func getAnalyzer(cliArgs CliArgs) Analyzer {
	switch cliArgs.Type {
	case "exitcode":
		// todo: maybe put this logic into a factory function NewExitCodeAnalyzer
		// assert that Args are well-formed
		s := strings.Fields(cliArgs.Args)
		Category := s[0]
		Subcategory := s[1]
		return ExitCodeAnalyzer{Category: Category, Subcategory: Subcategory}
	case "sq_apply_diffs":
		return SqApplyDiffsAnalyzer{}
	default:
		panic(fmt.Sprintf("unknown analyzer type %v", cliArgs.Type))
	}

}

func main() {

	var (
		Type      string
		Args      string
		OutputDir string
	)

	flag.StringVar(&Type, "type", "", "analyzer type")
	flag.StringVar(&Args, "args", "", "args for the specified analyzer")
	flag.StringVar(&OutputDir, "output-dir", "artifacts/analysis", "directory to store analysis output")

	flag.Parse()

	cliArgs := CliArgs{Type: Type, Args: Args, OutputDir: OutputDir}

	// get analyzer
	analyzer := getAnalyzer(cliArgs)

	// run command, store stdout / stderr to log in buffer while printing to stdout
	var out bytes.Buffer
	mwriter := io.MultiWriter(&out, os.Stdout)
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmd.Stderr = mwriter
	cmd.Stdout = mwriter
	cmd.Run()
	code := cmd.ProcessState.ExitCode()

	// call analyzer with exitCode, log
	result := analyzer.analyze(AnalyzerArgs{Exitcode: code, Log: out.String()})

	resultB, _ := json.Marshal(result)
	fmt.Println(string(resultB))

	// cmd := exec.Command("bash", "-c", "true && true && echo \"yoyoyo\"")
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("output was: %s", out.String())
}
