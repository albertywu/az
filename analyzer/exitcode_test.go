package analyzer

import "testing"

func TestExitcodeAnalyzer(t *testing.T) {
	az := ExitCodeAnalyzer{Category: "infra", Subcategory: "checkout"}
	args := Args{
		Exitcode:        0,
		StdoutStr:       "",
		StderrStr:       "",
		StdoutStderrStr: "",
	}
	result := az.Run(args)
	if result.Category != "success" || result.Subcategory != "" {
		t.Errorf("expected result to be success")
	}
}
