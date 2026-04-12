package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *CommandService) Create(ctx context.Context, name string) (ws *domain.Workspace, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspace *domain.Workspace
	workspace, err = domain.NewWorkspace(sub.ID, name)
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

	var created *domain.Workspace
	created, err = uc.workspaces.Create(ctx, *workspace)
	if err != nil {
		return nil, err
	}

	return created, nil
}
