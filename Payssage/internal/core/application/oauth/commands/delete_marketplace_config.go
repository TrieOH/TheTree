package commands

import (
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *CommandService) DeleteMarketplaceConfig(ctx context.Context, workspaceName string) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DeleteMarketplaceConfig")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	return uc.marketplace.Delete(ctx, workspace.ID)
}
