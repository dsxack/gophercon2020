package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

type generateTask interface {
	Generate(file *ast.File) error
}

func main() {
	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatal("GOFILE env variable must be set")
	}

	astInFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}

	i := inspector.New([]*ast.File{astInFile})
	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	var tasks []generateTask

	i.Nodes(iFilter, func(node ast.Node, push bool) (proceed bool) {
		genDecl := node.(*ast.GenDecl)
		if genDecl.Doc == nil {
			return false
		}

		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return false
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}

		for _, comment := range genDecl.Doc.List {
			switch comment.Text {
			case "//repogen:entity":
				tasks = append(tasks, repositoryGenerator{
					typeSpec:   typeSpec,
					structType: structType,
				})
			}
		}

		return false
	})

	astOutFile := &ast.File{
		Name: astInFile.Name,
	}

	for _, g := range tasks {
		err = g.Generate(astOutFile)
		if err != nil {
			log.Fatalf("generate: %v", err)
		}
	}

	outFile, err := os.Create(strings.TrimSuffix(path, ".go") + "_gen.go")
	if err != nil {
		log.Fatalf("create file: %v", err)
	}

	err = printer.Fprint(outFile, token.NewFileSet(), astOutFile)
	if err != nil {
		log.Fatalf("print file: %v", err)
	}
}
