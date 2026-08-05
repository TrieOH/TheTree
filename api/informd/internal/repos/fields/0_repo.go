package fields

import (
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.FieldsRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("field"),
	}
}

func mapField(src sqlc.Field) models.Field {
	return models.Field{
		ID:           src.ID,
		StepID:       src.StepID,
		Key:          src.Key,
		Title:        src.Title,
		Description:  src.Description,
		PositionHint: src.PositionHint,
		Required:     src.Required,
		Type:         models.FieldType(src.Type),
		Placeholder:  src.Placeholder,
		DefaultValue: src.DefaultValue,
		Config:       src.Config,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func mapFieldSelectConfig(src sqlc.FieldSelectConfig) models.FieldSelectConfig {
	return models.FieldSelectConfig{
		FieldID:   src.FieldID,
		Behaviour: models.SelectBehaviour(src.Behaviour),
		ValueType: models.SelectValueType(src.ValueType),
		Options:   src.Options,
	}
}
