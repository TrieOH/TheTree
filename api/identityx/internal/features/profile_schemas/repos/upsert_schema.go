package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (r *schemaRepo) Upsert(ctx context.Context, schema models.ProjectProfileSchema) (*models.ProjectProfileSchema, error) {
	ctx, span := database.Span(ctx, r.tracer, "UpsertProfileSchema")
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
