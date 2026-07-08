package checkpoints

import (
	"context"

	"lib/database"
	"univents/contracts"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	checkpoints ports.CheckpointsRepository
	editions    ports.EditionsRepository
	logger      *zap.Logger
	tracer      trace.Tracer
	tx          database.TxRunner
}

func NewQueryService(
	checkpoints ports.CheckpointsRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		checkpoints: checkpoints,
		editions:    editions,
		logger:      logger,
		tracer:      tracer,
		tx:          tx,
	}
}

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []contracts.Checkpoint, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "CheckpointService.AdminList")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	return uc.checkpoints.List(ctx, edition.ID)
}
