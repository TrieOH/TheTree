package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ProfileSchemaRepo = (*Repo)(nil)

func NewSchemaRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) *Repo {
	return &Repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("project_profile_schema"),
	}
}

func mapProfileSchema(src sqlc.ProjectProfileSchema) models.ProjectProfileSchema {
	return models.ProjectProfileSchema{
		ProjectID: src.ProjectID,
		Schema:    src.Schema,
		Version:   src.Version,
		Active:    src.Active,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}
