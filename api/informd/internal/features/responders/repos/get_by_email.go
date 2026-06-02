package repos

import (
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) GetByEmail(ctx context.Context, email string) (*models.Responder, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponderRepo.GetByEmail")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetResponderByEmail(ctx, email)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponder(row)), nil
}
