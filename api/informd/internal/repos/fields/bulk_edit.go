package fields

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
)

func (repo *Repo) BulkEdit(ctx context.Context, fields []models.Field) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.BulkEdit")
	defer span.End()
	params := xslices.MapSlice(fields, ToBulkEditFieldsParams)
	return database.BatchExec(
		database.Queries(ctx, repo.q).BulkEditFields(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].ID.String() },
	)
}

func ToBulkEditFieldsParams(f models.Field) sqlc.BulkEditFieldsParams {
	return sqlc.BulkEditFieldsParams{
		ID:           f.ID,
		StepID:       f.StepID,
		Key:          f.Key,
		Title:        f.Title,
		Description:  f.Description,
		PositionHint: f.PositionHint,
		Required:     f.Required,
		Type:         string(f.Type),
		Placeholder:  f.Placeholder,
		DefaultValue: f.DefaultValue,
		Config:       f.Config,
	}
}
