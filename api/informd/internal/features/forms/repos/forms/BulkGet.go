package forms

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *formRepo) BulkGet(ctx context.Context, ids []uuid.UUID, params models.BulkGetParams) ([]models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "BulkGet")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).BulkGetForms(ctx, ids)
	if err != nil {
		return nil, repo.dbe.DB(err, "form")
	}
	forms := xslices.MapSlice(sqlcForms, mapForm)
	forms, err = models.FilterForms(forms, params)
	if err != nil {
		return nil, err
	}
	models.SortForms(forms, params)
	return forms, nil
}
