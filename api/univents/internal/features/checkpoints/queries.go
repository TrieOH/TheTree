package checkpoints

import (
	"context"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	checkpoints ports.CheckpointsRepository
	editions    ports.EditionsRepository
	tracer      trace.Tracer
	az          *authzed.Client
	tx          database.TxRunner
}

func NewQueryService(
	checkpoints ports.CheckpointsRepository,
	editions ports.EditionsRepository,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		checkpoints: checkpoints,
		editions:    editions,
		tracer:      tracer,
		az:          az,
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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_checkpoints"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.checkpoints.List(ctx, editionID)
}
