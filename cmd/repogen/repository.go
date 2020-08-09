package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var repositoryTemplate = template.Must(template.New("").Parse(`
package main

import (
    "github.com/jinzhu/gorm"
)

type {{ .EntityName }}Repository struct {
	db *gorm.DB
}

func New{{ .EntityName }}Repository(db *gorm.DB) {{ .EntityName }}Repository {
	return {{ .EntityName }}Repository{ db: db }
}

func (r {{ .EntityName }}Repository) Get({{ .PrimaryName }} {{ .PrimaryType }}) (*{{ .EntityName }}, error) {
entity := new({{ .EntityName }})
	err := r.db.Limit(1).Where("{{ .PrimarySQLName }} = ?", {{ .PrimaryName }}).Find(entity).Error
return entity, err
}

func (r {{ .EntityName }}Repository) Create(entity *{{ .EntityName }}) error {
	return r.db.Create(entity).Error
}

func (r {{ .EntityName }}Repository) Update(entity *{{ .EntityName }}) error {
	return r.db.Model(entity).Update(entity).Error
}

func (r {{ .EntityName }}Repository) Delete(entity *{{ .EntityName }}) error {
	return r.db.Delete(entity).Error
}
`))

type repositoryGenerator struct {
	typeSpec   *ast.TypeSpec
	structType *ast.StructType
}

func (r repositoryGenerator) Generate(file *ast.File) error {
	primary, err := r.primaryField()
	if err != nil {
		return err
	}

	type templateParams struct {
		EntityName     string
		PrimaryType    string
		PrimaryName    string
		PrimarySQLName string
	}

	params := templateParams{
		EntityName:     r.typeSpec.Name.Name,
		PrimaryName:    strcase.ToLowerCamel(primary.Names[0].Name),
		PrimarySQLName: strcase.ToSnake(primary.Names[0].Name),
		PrimaryType:    expr2string(primary.Type),
	}

	var buf bytes.Buffer
	err = repositoryTemplate.Execute(&buf, params)
	if err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	templateAst, err := parser.ParseFile(token.NewFileSet(), "", buf.Bytes(), parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse template: %v", err)
	}

	for _, decl := range templateAst.Decls {
		file.Decls = append(file.Decls, decl)
	}

	return nil
}

func (r repositoryGenerator) primaryField() (*ast.Field, error) {
	for _, field := range r.structType.Fields.List {
		if !strings.Contains(field.Tag.Value, "primary") {
			continue
		}

		return field, nil
	}

	return nil, fmt.Errorf("has no primary field")
}
