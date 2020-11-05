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
)

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
