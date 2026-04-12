package commands

import (
	"TriePayments/internal/shared/authz"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) RevokeAPIKey(ctx context.Context, workspaceName string, keyID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RevokeAPIKey")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("revoke_api_keys"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return err
	}

	if _, err := uc.apiKeys.Revoke(ctx, keyID, workspace.ID); err != nil {
		return err
	}

	return nil
}
