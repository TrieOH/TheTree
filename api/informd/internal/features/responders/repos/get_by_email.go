package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) GetByEmail(ctx context.Context, email string) (*models.Responder, error) {
	ctx, span := telemetry.StartSpan(ctx, "ResponderRepo.GetByEmail")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetResponderByEmail(ctx, email)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponder(row)), nil
}
