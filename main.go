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

	analyzer, err := getAnalyzer(opts)
	if err != nil {
		log.Fatalf("invalid analyzer type %v", opts.id)
	}

	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	var stdoutBuf, stderrBuf, combinedBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf, &combinedBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf, &combinedBuf)

	cmd.Run()
	stdoutStr, stderrStr, stdoutStderrStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()), string(combinedBuf.Bytes())

	code := cmd.ProcessState.ExitCode()
	result := analyzer.run(
		args{
			exitcode:        code,
			stdoutStr:       stdoutStr,
			stderrStr:       stderrStr,
			stdoutStderrStr: stdoutStderrStr,
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
		[]byte(fmt.Sprintf("%s %s", result.Category, result.Subcategory)),
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
