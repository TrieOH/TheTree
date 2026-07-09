package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) Finish(ctx context.Context, activityID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Finish")
	defer span.End()

	err := database.Queries(ctx, repo.q).FinishActivity(ctx, activityID)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
