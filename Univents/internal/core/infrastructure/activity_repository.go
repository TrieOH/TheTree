package infrastructure

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"

	"github.com/MintzyG/fail/v3"
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
		StartsAt:          toCreate.CreatedAt,
		EndsAt:            toCreate.CreatedAt,
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
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return mapActivityFromDB(&sqlcActivity), nil
}
