package repos

import (
	"lib/database"

	"univents/contracts"
	"univents/internal/database/sqlc"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type activitiesRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ActivitiesRepository = (*activitiesRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ActivitiesRepository {
	return &activitiesRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("activities"),
	}
}

func mapActivityFromDB(src *sqlc.Activity) *contracts.Activity {
	return &contracts.Activity{
		ID:                src.ID,
		EditionID:         src.EditionID,
		Title:             src.Title,
		Description:       src.Description,
		Status:            contracts.ActivityStatus(src.Status),
		Location:          src.Location,
		StartsAt:          src.StartsAt,
		EndsAt:            src.EndsAt,
		PresenterName:     src.PresenterName,
		TokenCost:         src.TokenCost,
		HasCapacity:       src.HasCapacity,
		Capacity:          src.Capacity,
		RemainingCapacity: src.RemainingCapacity,
		Difficulty:        src.Difficulty,
		CreatedBy:         src.CreatedBy,
		CreatedAt:         src.CreatedAt,
		UpdatedAt:         src.UpdatedAt,
		DeletedAt:         src.DeletedAt,
	}
}

func mapAttendanceRecordFromDB(src *sqlc.AttendanceRecord) *contracts.AttendanceRecord {
	return &contracts.AttendanceRecord{
		ID:          src.ID,
		ActivityID:  src.ActivityID,
		UserID:      src.UserID,
		Status:      contracts.AttendanceStatus(src.Status),
		CheckedInAt: src.CheckedInAt,
		CancelledAt: src.CancelledAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}
