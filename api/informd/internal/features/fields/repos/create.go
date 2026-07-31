package repos

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Field) (*models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.Create")
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
