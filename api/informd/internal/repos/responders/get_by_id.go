package responders

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Responder, error) {
	ctx, span := telemetry.StartSpan(ctx, "ResponderRepo.GetByID")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetResponderByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponder(row)), nil
}
