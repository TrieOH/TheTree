package checkpoints

import (
	"context"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	checkpoints ports.CheckpointsRepository
	editions    ports.EditionsRepository
	logger      *zap.Logger
	tracer      trace.Tracer
	tx          database.TxRunner
}

func NewCommandService(
	checkpoints ports.CheckpointsRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		checkpoints: checkpoints,
		editions:    editions,
		logger:      logger,
		tracer:      tracer,
		tx:          tx,
	}
}

func (uc *CommandService) Create(ctx context.Context, in contracts.CreateCheckpointSpec) (out *contracts.Checkpoint, err error) {
	ctx, span := uc.tracer.Start(ctx, "CheckpointService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	var validCheckpoint *contracts.Checkpoint
	validCheckpoint, err = contracts.NewCheckpoint(sub.ID, in, edition)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_checkpoints"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	var created *contracts.Checkpoint
	created, err = uc.checkpoints.Create(ctx, validCheckpoint)
	if err != nil {
		return nil, err
	}

	return created, nil
}
