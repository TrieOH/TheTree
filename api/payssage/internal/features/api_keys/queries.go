package api_keys

import (
	"context"

	"payssage/internal/platform/database"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	apiKeys    ports.ApiKeysRepo
	workspaces ports.WorkspaceRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueryService(
	apiKeys ports.ApiKeysRepo,
	workspaces ports.WorkspaceRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		apiKeys:    apiKeys,
		workspaces: workspaces,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *QueryService) List(ctx context.Context, workspaceName string) (ak []contracts.APIKey, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspace *contracts.Workspace
	workspace, err = uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	var keys []contracts.APIKey
	keys, err = uc.apiKeys.ListByWorkspace(ctx, workspace.ID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
