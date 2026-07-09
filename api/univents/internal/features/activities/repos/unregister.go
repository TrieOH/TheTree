package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) Unregister(ctx context.Context, userID, activityID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Unregister")
	defer span.End()

	err := database.Queries(ctx, repo.q).UnregisterFromActivity(ctx, sqlc.UnregisterFromActivityParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return repo.dbe(err, "attendance record")
	}

	return nil
}
