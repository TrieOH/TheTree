package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProfileSchemaRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrProfileSchemaNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (r *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, r.q) }

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

func (r *Repo) Get(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	result, err := r.db(ctx).GetProfileSchema(ctx, projectID)
	return database.One(result, err, r.dbe, mapProfileSchema)
}

func (r *Repo) Upsert(ctx context.Context, schema models.ProjectProfileSchema) (*models.ProjectProfileSchema, error) {
	result, err := r.db(ctx).UpsertProfileSchema(ctx, sqlc.UpsertProfileSchemaParams{
		ProjectID: schema.ProjectID,
		Schema:    schema.Schema,
		Active:    schema.Active,
	})
	return database.One(result, err, r.dbe, mapUpsertedProfileSchema)
}
