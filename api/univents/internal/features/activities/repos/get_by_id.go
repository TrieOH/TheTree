package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetByID")
	defer span.End()

	sqlcActivity, err := database.Queries(ctx, repo.q).GetActivityByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapActivityFromDB(&sqlcActivity), nil
}
