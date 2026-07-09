package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *activitiesRepo) Register(ctx context.Context, toRegister contracts.AttendanceRecord) (*contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Register")
	defer span.End()

	sqlcRecord, err := database.Queries(ctx, repo.q).RegisterToActivity(ctx, sqlc.RegisterToActivityParams{
		ActivityID: toRegister.ActivityID,
		UserID:     toRegister.UserID,
		Status:     sqlc.AttendanceStatus(toRegister.Status),
	})
	if err != nil {
		return nil, repo.dbe(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcRecord), nil
}
