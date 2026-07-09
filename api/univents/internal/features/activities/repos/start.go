package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) Start(ctx context.Context, activityID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Start")
	defer span.End()

	err := database.Queries(ctx, repo.q).StartActivity(ctx, activityID)
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
