package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
)

func expr2string(expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), expr)
	if err != nil {
		log.Fatalf("error print expression to string: %v", err)
	}
	return buf.String()
}


