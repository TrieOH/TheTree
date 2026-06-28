package api_keys

import (
	"context"
	"payssage/models"
	"payssage/ports"

	"lib/authz"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	apiKeys    ports.ApiKeysRepo
	workspaces ports.WorkspaceRepo
	logger     *zap.Logger
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueryService(
	apiKeys ports.ApiKeysRepo,
	workspaces ports.WorkspaceRepo,
	logger *zap.Logger,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		apiKeys:    apiKeys,
		workspaces: workspaces,
		logger:     logger,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *QueryService) List(ctx context.Context, workspaceName string) (ak []models.APIKey, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspace *models.Workspace
	workspace, err = uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	var keys []models.APIKey
	keys, err = uc.apiKeys.ListByWorkspace(ctx, workspace.ID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
