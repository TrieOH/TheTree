package activities

import (
	"context"

	"univents/internal/platform/database"
	"univents/internal/platform/database/sqlc"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type activitiesRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.ActivitiesRepository = (*activitiesRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ActivitiesRepository {
	return &activitiesRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *activitiesRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapActivityFromDB(src *sqlc.Activity) *contracts.Activity {
	return &contracts.Activity{
		ID:                src.ID,
		ScopeID:           src.ScopeID,
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

func (repo *activitiesRepo) Create(ctx context.Context, toCreate *contracts.Activity) (*contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Create")
	defer span.End()

	sqlcActivity, err := repo.queries(ctx).CreateActivity(ctx, sqlc.CreateActivityParams{
		ID:                toCreate.ID,
		EditionID:         toCreate.EditionID,
		Title:             toCreate.Title,
		Description:       toCreate.Description,
		Location:          toCreate.Location,
		StartsAt:          toCreate.StartsAt,
		EndsAt:            toCreate.EndsAt,
		PresenterName:     toCreate.PresenterName,
		TokenCost:         toCreate.TokenCost,
		HasCapacity:       toCreate.HasCapacity,
		Capacity:          toCreate.Capacity,
		RemainingCapacity: toCreate.RemainingCapacity,
		Difficulty:        toCreate.Difficulty,
		CreatedBy:         toCreate.CreatedBy,
		ScopeID:           toCreate.ScopeID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "activity")
	}

	return mapActivityFromDB(&sqlcActivity), nil
}

func (repo *activitiesRepo) Publish(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Publish")
	defer span.End()

	err := repo.queries(ctx).PublishActivity(ctx, id)
	if err != nil {
		return errx.FromDB(err, "activity")
	}

	return nil
}

func (repo *activitiesRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetByID")
	defer span.End()

	sqlcActivity, err := repo.queries(ctx).GetActivityByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "activity")
	}

	return mapActivityFromDB(&sqlcActivity), nil
}

func (repo *activitiesRepo) Start(ctx context.Context, activityID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Start")
	defer span.End()

	err := repo.queries(ctx).StartActivity(ctx, activityID)
	if err != nil {
		return errx.FromDB(err, "activity")
	}

	return nil
}

func (repo *activitiesRepo) Finish(ctx context.Context, activityID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Finish")
	defer span.End()

	err := repo.queries(ctx).FinishActivity(ctx, activityID)
	if err != nil {
		return errx.FromDB(err, "activity")
	}

	return nil
}

func (repo *activitiesRepo) List(ctx context.Context, editionID uuid.UUID) ([]contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.List")
	defer span.End()

	sqlcActivities, err := repo.queries(ctx).ListEditionActivities(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "activity")
	}

	out := make([]contracts.Activity, 0, len(sqlcActivities))
	for _, activity := range sqlcActivities {
		out = append(out, *mapActivityFromDB(&activity))
	}
	return out, nil
}

func (repo *activitiesRepo) ListAdmin(ctx context.Context, editionID uuid.UUID) ([]contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.ListAdmin")
	defer span.End()

	sqlcActivities, err := repo.queries(ctx).ListEditionActivitiesAdmin(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "activity")
	}

	out := make([]contracts.Activity, 0, len(sqlcActivities))
	for _, activity := range sqlcActivities {
		out = append(out, *mapActivityFromDB(&activity))
	}
	return out, nil
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

func (repo *activitiesRepo) Register(ctx context.Context, toRegister contracts.AttendanceRecord) (*contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Register")
	defer span.End()

	sqlcRecord, err := repo.queries(ctx).RegisterToActivity(ctx, sqlc.RegisterToActivityParams{
		ActivityID: toRegister.ActivityID,
		UserID:     toRegister.UserID,
		Status:     sqlc.AttendanceStatus(toRegister.Status),
	})
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcRecord), nil
}

func (repo *activitiesRepo) Unregister(ctx context.Context, userID, activityID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Unregister")
	defer span.End()

	err := repo.queries(ctx).UnregisterFromActivity(ctx, sqlc.UnregisterFromActivityParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return errx.FromDB(err, "attendance record")
	}

	return nil
}

func (repo *activitiesRepo) MarkAttendanceRecordStatus(ctx context.Context, id uuid.UUID, scannedBy *uuid.UUID, status contracts.AttendanceStatus) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.MarkAttendanceRecordStatus")
	defer span.End()

	err := repo.queries(ctx).MarkAttendanceRecordStatus(ctx, sqlc.MarkAttendanceRecordStatusParams{
		Status:    sqlc.AttendanceStatus(status),
		ScannedBy: scannedBy,
		ID:        id,
	})
	if err != nil {
		return errx.FromDB(err, "attendance record")
	}

	return nil
}

func (repo *activitiesRepo) GetAttendanceRecordByID(ctx context.Context, id uuid.UUID) (*contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetAttendanceRecordByID")
	defer span.End()

	sqlcAttendanceRecord, err := repo.queries(ctx).GetAttendanceRecordByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcAttendanceRecord), nil
}

func (repo *activitiesRepo) ListActivityAttendanceRecords(ctx context.Context, activityID uuid.UUID) ([]contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.ListActivityAttendanceRecords")
	defer span.End()

	sqlcAttendanceRecords, err := repo.queries(ctx).ListActivityAttendanceRecords(ctx, activityID)
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	out := make([]contracts.AttendanceRecord, 0, len(sqlcAttendanceRecords))
	for _, record := range sqlcAttendanceRecords {
		out = append(out, *mapAttendanceRecordFromDB(&record))
	}
	return out, nil
}

func (repo *activitiesRepo) GetActiveUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) (*contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	sqlcAttendanceRecord, err := repo.queries(ctx).GetActiveUserActivityAttendanceRecords(ctx, sqlc.GetActiveUserActivityAttendanceRecordsParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcAttendanceRecord), nil
}

func (repo *activitiesRepo) GetUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) ([]contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	sqlcRecords, err := repo.queries(ctx).GetUserActivityAttendanceRecords(ctx, sqlc.GetUserActivityAttendanceRecordsParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	out := make([]contracts.AttendanceRecord, 0, len(sqlcRecords))
	for _, record := range sqlcRecords {
		out = append(out, *mapAttendanceRecordFromDB(&record))
	}
	return out, nil
}

func (repo *activitiesRepo) IsRegistered(ctx context.Context, userID, activityID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	isRegistered, err := repo.queries(ctx).IsUserRegistered(ctx, sqlc.IsUserRegisteredParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return false, errx.FromDB(err, "attendance record")
	}

	return isRegistered, nil
}
