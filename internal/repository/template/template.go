package template

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"scm.x5.ru/x5m/go-backend/template/internal/entities/templates"
)

// TemplateRepository encapsulates squirrel Querry Builder and DB.Connection
type TemplateRepository struct {
	db   squirrel.StatementBuilderType
	conn *sql.DB
}

// New creates a new TemplateRepository
func New(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{
		db:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db), // Use PostgreSQL placeholder style
		conn: db,
	}
}

// Create inserts a new template into the database
// Returns newly created entity in db, or nil, error in another case.
func (r *TemplateRepository) Create(ctx context.Context, template *templates.Template) (*templates.Template, error) {
	// TODO: add Logging inside repository.
	query := r.db.Insert("template_example").
		Columns("first_column", "second_column", "created_at", "updated_at").
		Values(template.FirstColumn, template.SecondColumn, template.CreatedAt, template.UpdatedAt).
		Suffix("RETURNING *")

	err := query.QueryRow().Scan(&template.ID, &template.FirstColumn, &template.SecondColumn, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return template, nil
}

// FindByID retrieves a template by its ID
func (r *TemplateRepository) FindByID(ctx context.Context, id uint64) (*templates.Template, error) {
	var template templates.Template
	query := r.db.Select("*").
		From("template_example").
		Where("id = ?", id)

	err := query.QueryRow().Scan(&template.ID, &template.FirstColumn, &template.SecondColumn, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update modifies an existing template
func (r *TemplateRepository) Update(ctx context.Context, template *templates.Template) (*templates.Template, error) {
	query := r.db.Update("template_example").
		Set("updated_at", time.Now()).
		Set("first_column", template.FirstColumn).
		Set("second_column", template.SecondColumn).
		Where("id = ?", template.ID).
		Suffix("RETURNING *")

	err := query.QueryRow().Scan(&template.ID, &template.FirstColumn, &template.SecondColumn, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return template, nil
}

// Delete removes a template from the database
func (r *TemplateRepository) Delete(ctx context.Context, id uint64) (uint64, error) {
	query := r.db.Delete("template_example").
		Where("id = ?", id)

	_, err := query.Exec()
	if err != nil {
		return 0, err
	}
	return id, nil
}
