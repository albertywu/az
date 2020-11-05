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

func main() {

	var (
		id        string
		config    string
		outputDir string
	)

	flag.StringVar(&id, "type", "", "analyzer type")
	flag.StringVar(&config, "config", "", "args for the specified analyzer")
	flag.StringVar(&outputDir, "output-dir", "artifacts/analysis", "directory to store analysis output")

	flag.Parse()

	opts := opts{id: id, config: config}

	// get analyzer
	analyzer, err := getAnalyzer(opts)
	if err != nil {
		log.Fatalf("invalid analyzer type %v", opts.id)
	}

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
	result := analyzer.run(args{exitcode: code, log: out.String()})
	resultB, _ := json.Marshal(result)
	fmt.Println(string(resultB))

}
