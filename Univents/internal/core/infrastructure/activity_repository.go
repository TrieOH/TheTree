package infrastructure

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/errx"

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

var _ domain.ActivitiesRepository = (*activitiesRepo)(nil)

func NewActivityRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.ActivitiesRepository {
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

func mapActivityFromDB(src *sqlc.Activity) *domain.Activity {
	return &domain.Activity{
		ID:                src.ID,
		ScopeID:           src.ScopeID,
		EditionID:         src.EditionID,
		Title:             src.Title,
		Description:       src.Description,
		Status:            domain.ActivityStatus(src.Status),
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

func (repo *activitiesRepo) Create(ctx context.Context, toCreate *domain.Activity) (*domain.Activity, error) {
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

func (repo *activitiesRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Activity, error) {
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

func (repo *activitiesRepo) List(ctx context.Context, editionID uuid.UUID) ([]domain.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.List")
	defer span.End()

	sqlcActivities, err := repo.queries(ctx).ListEditionActivities(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "activity")
	}

	out := make([]domain.Activity, 0, len(sqlcActivities))
	for _, activity := range sqlcActivities {
		out = append(out, *mapActivityFromDB(&activity))
	}
	return out, nil
}

func (repo *activitiesRepo) ListAdmin(ctx context.Context, editionID uuid.UUID) ([]domain.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.ListAdmin")
	defer span.End()

	sqlcActivities, err := repo.queries(ctx).ListEditionActivitiesAdmin(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "activity")
	}

	out := make([]domain.Activity, 0, len(sqlcActivities))
	for _, activity := range sqlcActivities {
		out = append(out, *mapActivityFromDB(&activity))
	}
	return out, nil
}

func mapAttendanceRecordFromDB(src *sqlc.AttendanceRecord) *domain.AttendanceRecord {
	return &domain.AttendanceRecord{
		ID:          src.ID,
		ActivityID:  src.ActivityID,
		UserID:      src.UserID,
		Status:      domain.AttendanceStatus(src.Status),
		CheckedInAt: src.CheckedInAt,
		CancelledAt: src.CancelledAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func (repo *activitiesRepo) Register(ctx context.Context, toRegister domain.AttendanceRecord) (*domain.AttendanceRecord, error) {
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

func (repo *activitiesRepo) MarkAttendanceRecordStatus(ctx context.Context, id uuid.UUID, scannedBy *uuid.UUID, status domain.AttendanceStatus) error {
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

func (repo *activitiesRepo) GetAttendanceRecordByID(ctx context.Context, id uuid.UUID) (*domain.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetAttendanceRecordByID")
	defer span.End()

	sqlcAttendanceRecord, err := repo.queries(ctx).GetAttendanceRecordByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcAttendanceRecord), nil
}

func (repo *activitiesRepo) ListActivityAttendanceRecords(ctx context.Context, activityID uuid.UUID) ([]domain.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.ListActivityAttendanceRecords")
	defer span.End()

	sqlcAttendanceRecords, err := repo.queries(ctx).ListActivityAttendanceRecords(ctx, activityID)
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	out := make([]domain.AttendanceRecord, 0, len(sqlcAttendanceRecords))
	for _, record := range sqlcAttendanceRecords {
		out = append(out, *mapAttendanceRecordFromDB(&record))
	}
	return out, nil
}

func (repo *activitiesRepo) GetActiveUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) (*domain.AttendanceRecord, error) {
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

func (repo *activitiesRepo) GetUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) ([]domain.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	sqlcRecords, err := repo.queries(ctx).GetUserActivityAttendanceRecords(ctx, sqlc.GetUserActivityAttendanceRecordsParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "attendance record")
	}

	out := make([]domain.AttendanceRecord, 0, len(sqlcRecords))
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

func mapExportRowFromDB(src sqlc.AttendanceExportRow) domain.AttendanceExportRow {
	row := domain.AttendanceExportRow{
		AttendanceID:     src.AttendanceID,
		UserID:           src.UserID,
		AttendanceStatus: domain.AttendanceStatus(src.AttendanceStatus),
		CheckedInAt:      src.CheckedInAt,
		CancelledAt:      src.CancelledAt,
		ScannedBy:        src.ScannedBy,
		RegisteredAt:     src.CreatedAt,
		ActivityID:       src.ActivityID,
		ActivityTitle:    src.ActivityTitle,
		ActivityLocation: src.ActivityLocation,
		ActivityStatus:   domain.ActivityStatus(src.ActivityStatus),
		ActivityStartsAt: src.ActivityStartsAt,
		ActivityEndsAt:   src.ActivityEndsAt,
	}

	if src.Difficulty != nil {
		d := *src.Difficulty
		row.Difficulty = &d
	}

	return row
}
func filtersToParams(editionID uuid.UUID, f domain.ExportFilters) sqlc.AttendanceExportParams {
	return sqlc.AttendanceExportParams{
		EditionID:          editionID,
		ActivityIds:        toUUIDSlice(f.ActivityIDs),
		AttendanceStatuses: toSlice(f.AttendanceStatuses, func(s domain.AttendanceStatus) sqlc.AttendanceStatus { return sqlc.AttendanceStatus(s) }),
		ActivityStatuses:   toSlice(f.ActivityStatuses, func(s domain.ActivityStatus) sqlc.ActivityStatus { return sqlc.ActivityStatus(s) }),
		Difficulties:       f.Difficulties,
		DateFrom:           f.DateFrom,
		DateTo:             f.DateTo,
	}
}

func (repo *activitiesRepo) AttendanceExport(ctx context.Context, editionID uuid.UUID, filters domain.ExportFilters) ([]domain.AttendanceExportRow, error) {
	ctx, span := repo.tracer.Start(ctx, "AttendanceExportRepo.AttendanceExport")
	defer span.End()

	sqlcRows, err := repo.queries(ctx).AttendanceExport(ctx, filtersToParams(editionID, filters))
	if err != nil {
		telemetry.Log().Error("raw db error", zap.Error(err))
		return nil, errx.FromDB(err, "attendance_export")
	}

	rows := make([]domain.AttendanceExportRow, len(sqlcRows))
	for i, r := range sqlcRows {
		rows[i] = mapExportRowFromDB(r)
	}

	return rows, nil
}

// toUUIDSlice returns nil when empty so the SQL NULL check skips the filter.
func toUUIDSlice(ids []uuid.UUID) []uuid.UUID {
	if len(ids) == 0 {
		return nil
	}
	return ids
}

// toSlice maps a typed slice to another type, returning nil when empty.
func toSlice[In any, Out any](slice []In, fn func(In) Out) []Out {
	if len(slice) == 0 {
		return nil
	}
	out := make([]Out, len(slice))
	for i, v := range slice {
		out[i] = fn(v)
	}
	return out
}
