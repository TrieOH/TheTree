package repo

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Form) (*models.Form, error) {
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
