package checkpoints

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

type checkpointsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.CheckpointsRepository = (*checkpointsRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.CheckpointsRepository {
	return &checkpointsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *checkpointsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapCheckpointFromDB(src *sqlc.Checkpoint) *contracts.Checkpoint {
	return &contracts.Checkpoint{
		ID:         src.ID,
		ScopeID:    src.ScopeID,
		EditionID:  src.EditionID,
		Name:       src.Name,
		Type:       contracts.CheckpointType(src.Type),
		AccessMode: contracts.CheckpointAccess(src.AccessMode),
		StartsAt:   src.StartsAt,
		EndsAt:     src.EndsAt,
		CreatedBy:  src.CreatedBy,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
		DeletedAt:  src.DeletedAt,
	}
}

func (repo *checkpointsRepo) Create(ctx context.Context, toCreate *contracts.Checkpoint) (*contracts.Checkpoint, error) {
	ctx, span := repo.tracer.Start(ctx, "CheckpointRepo.Create")
	defer span.End()

	sqlcCheckpoint, err := repo.queries(ctx).CreateCheckpoint(ctx, sqlc.CreateCheckpointParams{
		ID:         toCreate.ID,
		ScopeID:    toCreate.ScopeID,
		EditionID:  toCreate.EditionID,
		Name:       toCreate.Name,
		Type:       sqlc.CheckpointType(toCreate.Type),
		AccessMode: sqlc.CheckpointAccess(toCreate.AccessMode),
		StartsAt:   toCreate.StartsAt,
		EndsAt:     toCreate.EndsAt,
		CreatedBy:  toCreate.CreatedBy,
	})
	if err != nil {
		return nil, errx.FromDB(err, "checkpoint")
	}

	return mapCheckpointFromDB(&sqlcCheckpoint), nil
}

func (repo *checkpointsRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Checkpoint, error) {
	ctx, span := repo.tracer.Start(ctx, "CheckpointRepo.GetByID")
	defer span.End()

	sqlcCheckpoint, err := repo.queries(ctx).GetCheckpointByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "checkpoint")
	}

	return mapCheckpointFromDB(&sqlcCheckpoint), nil
}

func (repo *checkpointsRepo) List(ctx context.Context, editionID uuid.UUID) ([]contracts.Checkpoint, error) {
	ctx, span := repo.tracer.Start(ctx, "CheckpointRepo.List")
	defer span.End()

	sqlcCheckpoints, err := repo.queries(ctx).ListEditionCheckpoints(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "checkpoint")
	}

	out := make([]contracts.Checkpoint, 0, len(sqlcCheckpoints))
	for _, checkpoint := range sqlcCheckpoints {
		out = append(out, *mapCheckpointFromDB(&checkpoint))
	}
	return out, nil
}
