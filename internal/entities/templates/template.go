package templates

import (
	"time"

	"scm.x5.ru/x5m/go-backend/template/internal/entities"
)

// Template is an example entity for describing DB model
type Template struct {
	entities.IDModel
	entities.TimeModel
	FirstColumn  string `db:"first_column"`
	SecondColumn string `db:"second_column"`
}

// New creates TemplateEntity wihtout assigning ID field, because it serial type
func New(first, second string) *Template {
	return &Template{
		FirstColumn:  first,
		SecondColumn: second,
		TimeModel: entities.TimeModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
