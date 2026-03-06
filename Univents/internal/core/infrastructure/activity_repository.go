package infrastructure

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
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

	out := make([]domain.Activity, len(sqlcActivities))
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

	out := make([]domain.Activity, len(sqlcActivities))
	for _, activity := range sqlcActivities {
		out = append(out, *mapActivityFromDB(&activity))
	}
	return out, nil
}
