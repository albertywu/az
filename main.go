package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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

func getAnalyzer(cliArgs CliArgs) (Analyzer, error) {
	switch cliArgs.Type {
	case "exitcode":
		s := strings.Fields(cliArgs.Args)
		Category := s[0]
		Subcategory := s[1]
		// TODO: assert that category is one of { success, canceled, infra_failure, user_failure }
		return ExitCodeAnalyzer{Category: Category, Subcategory: Subcategory}, nil
	case "sq_apply_diffs":
		return SqApplyDiffsAnalyzer{}, nil
	default:
		return nil, fmt.Errorf("invalid analyzer type %v", cliArgs.Type)
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
	analyzer, err := getAnalyzer(cliArgs)
	if err != nil {
		log.Fatalf("invalid analyzer type %v", cliArgs.Type)
	}

	// run command, store stdout / stderr to log in buffer while printing to stdout
	var out bytes.Buffer
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("could not get stderr pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("could not get stdout pipe: %v", err)
	}
	go func() {
		merged := io.MultiReader(stderr, stdout)
		scanner := bufio.NewScanner(merged)
		for scanner.Scan() {
			msg := scanner.Text()
			fmt.Println(msg)
			out.Write(scanner.Bytes())
		}
	}()

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
