package profile_schemas

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProfileSchemaRepo = (*Repo)(nil)

func NewSchemaRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("project_profile_schema"),
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

// mapUpsertedProfileSchema maps the versioned upsert's result row. The two
// sqlc types share an identical shape, so a plain conversion suffices.
func mapUpsertedProfileSchema(src sqlc.UpsertProfileSchemaRow) models.ProjectProfileSchema {
	return mapProfileSchema(sqlc.ProjectProfileSchema(src))
}
