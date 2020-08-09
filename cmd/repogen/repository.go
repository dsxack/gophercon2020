package main

import (
	"go/ast"
	"text/template"
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
	return nil
}
