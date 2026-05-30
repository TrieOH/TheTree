package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Field) (*models.Field, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.Create")
	defer span.End()
	sqlcField, err := database.Queries(ctx, repo.q).CreateField(ctx, sqlc.CreateFieldParams{
		StepID:       toCreate.StepID,
		Key:          toCreate.Key,
		Title:        toCreate.Title,
		Description:  toCreate.Description,
		PositionHint: toCreate.PositionHint,
		Required:     toCreate.Required,
		Type:         string(toCreate.Type),
		Placeholder:  toCreate.Placeholder,
		DefaultValue: toCreate.DefaultValue,
		Config:       toCreate.Config,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapField(sqlcField)), nil
}
