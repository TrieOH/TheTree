package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
)

func (uc *CommandService) DisableSandbox(ctx context.Context, workspaceName string) (*domain.Workspace, error) {
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
