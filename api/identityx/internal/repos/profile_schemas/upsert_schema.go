package profile_schemas

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (r *Repo) Upsert(ctx context.Context, schema models.ProjectProfileSchema) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "Upsert")
	defer span.End()

	result, err := database.Queries(ctx, r.q).UpsertProfileSchema(ctx, sqlc.UpsertProfileSchemaParams{
		ProjectID: schema.ProjectID,
		Schema:    schema.Schema,
		Active:    schema.Active,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapProfileSchema(result)), nil
}
