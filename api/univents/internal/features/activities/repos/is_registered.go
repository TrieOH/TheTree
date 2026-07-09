package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) IsRegistered(ctx context.Context, userID, activityID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	isRegistered, err := database.Queries(ctx, repo.q).IsUserRegistered(ctx, sqlc.IsUserRegisteredParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return false, repo.dbe(err, "attendance record")
	}

	return isRegistered, nil
}
