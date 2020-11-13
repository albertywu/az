package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type artifactTest struct {
	expectedOutput       string
	expectedStdout       string
	expectedStderr       string
	expectedStdoutStderr string
}

func checkArtifacts(t *testing.T, test artifactTest) {
	content, _ := ioutil.ReadFile("artifacts/analysis/output")
	if string(content) != test.expectedOutput {
		t.Errorf(fmt.Sprintf("expected artifacts/analysis/output to contain %s, but got %s", test.expectedOutput, content))
	}
	stdout, _ := ioutil.ReadFile("artifacts/analysis/stdout")
	if string(stdout) != test.expectedStdout {
		t.Errorf(fmt.Sprintf("expected to log %s to artifacts/analysis/stdout, but got: %s", test.expectedStdout, stdout))
	}
	stderr, _ := ioutil.ReadFile("artifacts/analysis/stderr")
	if string(stderr) != test.expectedStderr {
		t.Errorf(fmt.Sprintf("expected to log %s to artifacts/analysis/stderr, but got: %s", test.expectedStderr, stderr))
	}
	stdoutStderr, _ := ioutil.ReadFile("artifacts/analysis/stdoutStderr")
	if string(stdoutStderr) != test.expectedStdoutStderr {
		t.Errorf(fmt.Sprintf("expected to log %s to artifacts/analysis/stdoutStderr, but got: %s", test.expectedStdoutStderr, stdoutStderr))
	}
}

func makeBinary() {
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		"az",
	)
	cmd.Run()
}

func cleanup() {
	os.Remove("./az")
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

func TestExitcodeSuccess(t *testing.T) {
	config := "infra foo"
	cmd := exec.Command(
		"./az",
		"--type",
		"exitcode",
		"--config",
		config,
		"fixtures/exit_zero.sh",
	)
	cmd.Run()
	code := cmd.ProcessState.ExitCode()
	if code != 0 {
		t.Errorf("expected process to exit with code 0")
	}

	checkArtifacts(t, artifactTest{
		expectedOutput:       "success",
		expectedStdout:       "foo\n",
		expectedStderr:       "bar\n",
		expectedStdoutStderr: "foo\nbar\n",
	})
}

func TestExitcodeFailure(t *testing.T) {
	config := "infra foo"
	cmd := exec.Command(
		"./az",
		"--type",
		"exitcode",
		"--config",
		config,
		"fixtures/exit_nonzero.sh",
	)

	cmd.Run()
	code := cmd.ProcessState.ExitCode()
	if code != 1 {
		t.Errorf("expected process to exit with code 1")
	}

	checkArtifacts(t, artifactTest{
		expectedOutput:       config,
		expectedStdout:       "foo\n",
		expectedStderr:       "bar\n",
		expectedStdoutStderr: "foo\nbar\n",
	})
}
