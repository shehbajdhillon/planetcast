package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {

	flag.Parse()

	for _, filename := range flag.Args() {

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)

		if err != nil {
			fmt.Println("Failed to parse file:", err)
			os.Exit(1)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)

			if !ok {
				return true
			}

			if fn.Doc == nil {
				return true
			}

			for _, comment := range fn.Doc.List {
				if strings.HasPrefix(comment.Text, "// @PrcTrace") {
					fmt.Println(comment.Text)
					fmt.Println("Found function to instrument:", fn.Name.Name)
					injectTracing(fn)
				}
			}

			return true
		})

		buf := new(bytes.Buffer)
		if err := format.Node(buf, fset, node); err != nil {
			fmt.Printf("Could not format modified code: %v\n", err)
			continue
		}

		if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
			fmt.Printf("Could not write modified code to file: %v\n", err)
			continue
		}

	}
}

func injectTracing(fn *ast.FuncDecl) {
	hasContext := false

	for _, param := range fn.Type.Params.List {
		if expr, ok := param.Type.(*ast.SelectorExpr); ok {
			if pkg, ok := expr.X.(*ast.Ident); ok && pkg.Name == "context" && expr.Sel.Name == "Context" {
				hasContext = true
				break
			}
		}
	}

	if !hasContext {
		fmt.Println("Function does not have context passed in:", fn.Name.Name)
		return
	}

	ctxAssignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ctx"),
			ast.NewIdent("span"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("otel"),
							Sel: ast.NewIdent("Tracer"),
						},
						Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: `"tracer-name"`}},
					},
					Sel: ast.NewIdent("Start"),
				},
				Args: []ast.Expr{
					ast.NewIdent("ctx"),
					&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, fn.Name.Name)},
				},
			},
		},
	}

	deferStmt := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("span"),
				Sel: ast.NewIdent("End"),
			},
		},
	}
	fn.Body.List = append([]ast.Stmt{ctxAssignStmt, deferStmt}, fn.Body.List...)
}
