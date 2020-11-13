package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"./analyzer"
)

func main() {

	var (
		_type     string
		config    string
		outputDir string
	)

	flag.StringVar(&_type, "type", "", "analyzer type")
	flag.StringVar(&config, "config", "", "args for the specified analyzer")
	flag.StringVar(&outputDir, "output-dir", "artifacts/analysis", "directory to store analysis output")

	flag.Parse()

	opts := analyzer.Opts{Type: _type, Config: config}

	az, err := analyzer.GetAnalyzer(opts)
	if err != nil {
		log.Fatalf("invalid analyzer type %v", opts.Type)
	}

	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)

	var stdoutBuf, stderrBuf, combinedBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf, &combinedBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf, &combinedBuf)

	cmd.Run()
	stdoutStr, stderrStr, stdoutStderrStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()), string(combinedBuf.Bytes())

	code := cmd.ProcessState.ExitCode()

	result := az.Run(
		analyzer.Args{
			Exitcode:        code,
			StdoutStr:       stdoutStr,
			StderrStr:       stderrStr,
			StdoutStderrStr: stdoutStderrStr,
		},
	)

	var outputDirAbs string
	if path.IsAbs(outputDir) {
		outputDirAbs = outputDir
	} else {
		cwd, _ := os.Getwd()
		outputDirAbs = fmt.Sprintf("%s/%s", cwd, outputDir)
	}
	os.RemoveAll(outputDirAbs)
	err = os.MkdirAll(outputDirAbs, 0755)
	if err != nil {
		log.Fatalf("could not create dir at %v", outputDirAbs)
	}
	ioutil.WriteFile(
		fmt.Sprintf("%s/output", outputDirAbs),
		[]byte(
			strings.TrimSpace(
				fmt.Sprintf("%s %s", result.Category, result.Subcategory),
			),
		),
		0644,
	)
	ioutil.WriteFile(
		fmt.Sprintf("%s/stdout", outputDirAbs),
		[]byte(stdoutStr),
		0644,
	)
	ioutil.WriteFile(
		fmt.Sprintf("%s/stderr", outputDirAbs),
		[]byte(stderrStr),
		0644,
	)
	ioutil.WriteFile(
		fmt.Sprintf("%s/stdoutStderr", outputDirAbs),
		[]byte(stdoutStderrStr),
		0644,
	)

	os.Exit(code)
}
