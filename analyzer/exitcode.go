package analyzer

type ExitCodeAnalyzer struct {
	Category    string
	Subcategory string
}

func (a ExitCodeAnalyzer) Run(args Args) Result {
	if args.Exitcode == 0 {
		return Result{Category: "success", Subcategory: ""}
	}
	return Result{Category: a.Category, Subcategory: a.Subcategory}
}
