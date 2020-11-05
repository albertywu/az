package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
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
	result := analyzer.run(args{exitcode: code, log: out.String()})
	resultB, _ := json.Marshal(result)

	var failureDirAbs string
	if path.IsAbs(outputDir) {
		failureDirAbs = outputDir
	} else {
		cwd, _ := os.Getwd()
		failureDirAbs = fmt.Sprintf("%s/%s", cwd, outputDir)
	}
	os.RemoveAll(failureDirAbs)
	err = os.MkdirAll(failureDirAbs, 0755)
	if err != nil {
		log.Fatalf("could not create dir at %v", failureDirAbs)
	}
	ioutil.WriteFile(fmt.Sprintf("%s/failure", failureDirAbs), []byte(resultB), 0644)
}
