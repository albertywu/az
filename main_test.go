package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func makeBinary() {
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		"analyze",
	)
	cmd.Run()
}

func cleanup() {
	os.Remove("./analyze")
	os.RemoveAll("artifacts")
}

func TestMain(m *testing.M) {
	fmt.Println("...  creating analysis binary ...")
	makeBinary()
	code := m.Run()
	fmt.Println("...  removing analysis binary ...")
	cleanup()
	os.Exit(code)
}

func TestExitcodeAnalyzerPassing(t *testing.T) {
	cmd := exec.Command(
		"./analyze",
		"--type",
		"exitcode",
		"--config",
		"infra failure",
		"ls",
	)
	err := cmd.Run()
	if err != nil {
		t.Errorf("expected pass")
	}
	_, err = ioutil.ReadFile("artifacts/analysis/failure")
	if err == nil {
		t.Errorf("failure file should not exist")
	}
}

func TestExitcodeAnalyzerFailing(t *testing.T) {
	config := "infra checkout"
	cmd := exec.Command(
		"./analyze",
		"--type",
		"exitcode",
		"--config",
		config,
		"false",
	)
	err := cmd.Run()
	if err == nil {
		t.Errorf("expected failure")
	}
	fmt.Println("Exit code:", cmd.ProcessState.ExitCode())
	if cmd.ProcessState.ExitCode() != 1 {
		t.Errorf("expected process to exit with code 1")
	}
	content, _ := ioutil.ReadFile("artifacts/analysis/failure")
	if string(content) != config {
		t.Errorf(fmt.Sprintf("expected failure to be %s", config))
	}
}
