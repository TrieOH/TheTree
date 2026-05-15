package workspaces

import (
	"context"
	"payssage/internal/platform/database"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	workspaces ports.WorkspaceRepo
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommandService(
	workspaces ports.WorkspaceRepo,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		workspaces: workspaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *CommandService) Create(ctx context.Context, name string) (ws *contracts.Workspace, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspace *contracts.Workspace
	workspace, err = contracts.NewWorkspace(sub.ID, name)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_workspaces"),
		authz.Resource("platform", "global"),
	); err != nil {
		return nil, err
	}

	var created *contracts.Workspace
	created, err = uc.workspaces.Create(ctx, *workspace)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) DisableSandbox(ctx context.Context, workspaceName string) (*contracts.Workspace, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DisableSandbox")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			return nil, errx.NotFound("workspace")
		}
		return nil, err
	}

	return uc.workspaces.DisableSandbox(ctx, workspace.ID)
}

func (uc *CommandService) EnableSandbox(ctx context.Context, workspaceName string) (*contracts.Workspace, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.EnableSandbox")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			return nil, errx.NotFound("workspace")
		}
		return nil, err
	}

	return uc.workspaces.EnableSandbox(ctx, workspace.ID)
}
