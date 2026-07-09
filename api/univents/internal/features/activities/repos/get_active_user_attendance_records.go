package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) GetActiveUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) (*contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	sqlcAttendanceRecord, err := database.Queries(ctx, repo.q).GetActiveUserActivityAttendanceRecords(ctx, sqlc.GetActiveUserActivityAttendanceRecordsParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return nil, repo.dbe(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcAttendanceRecord), nil
}
