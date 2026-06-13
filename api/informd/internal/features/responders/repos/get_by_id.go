package repos

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Responder, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponderRepo.GetByID")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetResponderByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponder(row)), nil
}
