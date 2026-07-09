package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) ListAdmin(ctx context.Context, editionID uuid.UUID) ([]contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.ListAdmin")
	defer span.End()

	sqlcActivities, err := database.Queries(ctx, repo.q).ListEditionActivitiesAdmin(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	out := make([]contracts.Activity, 0, len(sqlcActivities))
	for _, activity := range sqlcActivities {
		out = append(out, *mapActivityFromDB(&activity))
	}
	return out, nil
}
