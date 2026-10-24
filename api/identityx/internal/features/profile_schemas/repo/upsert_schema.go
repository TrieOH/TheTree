package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (r *Repo) Upsert(ctx context.Context, schema models.ProjectProfileSchema) (*models.ProjectProfileSchema, error) {
	result, err := database.Queries(ctx, r.q).UpsertProfileSchema(ctx, sqlc.UpsertProfileSchemaParams{
		ProjectID: schema.ProjectID,
		Schema:    schema.Schema,
		Active:    schema.Active,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapUpsertedProfileSchema(result)), nil
}
