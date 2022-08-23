package main

import (
	"go/ast"
	"go/token"

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
	funcexpr := func(x *ast.FuncDecl) (exitCalled bool) {
		if x.Name.Name != "main" {
			return false
		}
		for _, stmt := range x.Body.List {
			var isExit, isOs bool
			if expr, ok := stmt.(*ast.ExprStmt); ok {
				if cExpr, ok := expr.X.(*ast.CallExpr); ok {
					if fun, ok := cExpr.Fun.(*ast.SelectorExpr); ok {
						if x, ok := fun.X.(*ast.Ident); ok {
							if x.Name == "os" {
								isOs = true
							}
						}
						if fun.Sel.Name == "Exit" {
							isExit = true
						}
						if isOs && isExit {
							return true
						}
					}
				}
			}
		}
		return false
	}

	pkgFunc := func(x *ast.Ident) (isMain bool) {
		return x.Name == "main"
	}

	for _, file := range pass.Files {
		var isMain, exitCalled bool
		var exitPos token.Pos
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.Ident:
				isMain = pkgFunc(x)
			case *ast.FuncDecl:
				exitCalled = funcexpr(x)
			}
			if isMain && exitCalled {
				pass.Reportf(exitPos, "os.Exit called")
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
		MainExitAnalyzer,
	}
	myChecks = append(myChecks, passesChecks...)

	multichecker.Main(
		myChecks...,
	)
}
