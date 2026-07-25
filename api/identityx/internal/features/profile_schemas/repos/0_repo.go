package repos

import (
	sqlc2 "IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.ProfileSchemaRepo = (*Repo)(nil)

func NewSchemaRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("project_profile_schema"),
	}
}

func mapProfileSchema(src sqlc2.ProjectProfileSchema) models.ProjectProfileSchema {
	return models.ProjectProfileSchema{
		ProjectID: src.ProjectID,
		Schema:    src.Schema,
		Version:   src.Version,
		Active:    src.Active,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}
