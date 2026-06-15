package workspaces

import (
	"context"
	"payssage/ports"

	"lib/authz"
	"lib/database"
	"payssage/models"

	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	workspaces ports.WorkspaceRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueryService(
	workspaces ports.WorkspaceRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		workspaces: workspaces,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *QueryService) List(ctx context.Context) (ws []models.Workspace, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspaces []models.Workspace
	workspaces, err = uc.workspaces.List(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}
