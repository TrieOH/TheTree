package forms

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Form) (*models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.Create")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).CreateForm(ctx, sqlc.CreateFormParams{
		NamespaceID: toCreate.NamespaceID,
		CreatedBy:   toCreate.CreatedBy,
		OwnerID:     toCreate.OwnerID,
		Name:        toCreate.Title,
		Status:      string(toCreate.Status),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
