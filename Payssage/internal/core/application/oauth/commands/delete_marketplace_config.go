package commands

import (
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) DeleteMarketplaceConfig(ctx context.Context, workspaceName string, credentialID uuid.UUID) error {
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

	// verify the credential belongs to this workspace before deleting
	cred, err := uc.credentials.GetByID(ctx, credentialID)
	if err != nil {
		return err
	}
	if cred.WorkspaceID != workspace.ID {
		return errx.Forbidden("credential").SetMessage("credential does not belong to this workspace")
	}

	return uc.marketplace.Delete(ctx, workspace.ID, credentialID)
}
