package responders

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Responder) (*models.Responder, error) {
	ctx, span := telemetry.StartSpan(ctx, "ResponderRepo.Create")
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
