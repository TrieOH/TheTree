package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionsRepo.GetByID")
	defer span.End()
	edition, err := database.Queries(ctx, repo.q).GetEditionByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(edition)), nil
}
