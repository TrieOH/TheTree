package workspaces

import (
	"context"
	"payssage/internal/platform/database"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/ports"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	workspaces ports.WorkspaceRepo
	gaClient   *goauth.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueryService(
	workspaces ports.WorkspaceRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		workspaces: workspaces,
		gaClient:   gaClient,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *QueryService) List(ctx context.Context) (ws []contracts.Workspace, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspaces []contracts.Workspace
	workspaces, err = uc.workspaces.List(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}
