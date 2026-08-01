package editions

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionsRepo.GetActive")
	defer span.End()
	edition, err := database.Queries(ctx, repo.q).GetActiveEdition(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(edition)), nil
}
