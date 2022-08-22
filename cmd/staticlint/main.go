package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func run(pass *analysis.Pass) (interface{}, error) {
	funcexpr := func(x *ast.FuncDecl) {
		for _, stmt := range x.Body.List {
			if expr, ok := stmt.(*ast.CallExpr); ok {
				if fun, ok := expr.Fun.(*ast.SelectorExpr); ok {
					if fun.Sel.Name == "os.Exit" {
						pass.Reportf(fun.Sel.NamePos, "using os.Exit in main")
					}
				}
			}
		}
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				funcexpr(x)
			}
			return true
		})
	}
	return nil, nil
}

var MainExitAnalyzer = &analysis.Analyzer{
	Name: "mainExit",
	Doc:  "check for os.Exit in main function",
	Run:  run,
}

var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}

func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named:
		return []bool{isErrorType(t)}
	case *types.Pointer:
		return []bool{isErrorType(t)}
	case *types.Tuple:
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}

func main() {
	var myChecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	for _, v := range quickfix.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	for _, v := range stylecheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	passesChecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
	}
	myChecks = append(myChecks, passesChecks...)

	multichecker.Main(
		myChecks...,
	)
}
