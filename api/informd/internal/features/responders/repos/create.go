package repos

import (
	"context"

	"Informd/internal/database/sqlc"
	"Informd/models"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Responder) (*models.Responder, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponderRepo.Create")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).CreateResponder(ctx, sqlc.CreateResponderParams{
		UserID: toCreate.UserID,
		Email:  toCreate.Email,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponder(row)), nil
}
