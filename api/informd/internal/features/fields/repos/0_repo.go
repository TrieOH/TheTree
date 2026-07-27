package repos

import (
	sqlc2 "Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.FieldsRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.FieldsRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("field"),
	}
}

func mapField(src sqlc2.Field) models.Field {
	return models.Field{
		ID:           src.ID,
		StepID:       src.StepID,
		Key:          src.Key,
		Title:        src.Title,
		Description:  src.Description,
		PositionHint: int(src.PositionHint),
		Required:     src.Required,
		Type:         models.FieldType(src.Type),
		Placeholder:  src.Placeholder,
		DefaultValue: src.DefaultValue,
		Config:       src.Config,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func mapFieldSelectConfig(src sqlc2.FieldSelectConfig) models.FieldSelectConfig {
	return models.FieldSelectConfig{
		FieldID:   src.FieldID,
		Behaviour: models.SelectBehaviour(src.Behaviour),
		ValueType: models.SelectValueType(src.ValueType),
		Options:   src.Options,
	}
}
