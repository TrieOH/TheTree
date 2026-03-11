package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) RevokeCredential(ctx context.Context, workspaceName string, credentialID uuid.UUID) (*domain.ProviderCredential, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RevokeCredential")
	defer span.End()

	subject, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, subject.ID)
	if err != nil {
		return nil, err
	}

	cred, err := uc.credentials.Revoke(ctx, credentialID, workspace.ID)
	if err != nil {
		return nil, err
	}

	// if this credential was backing a marketplace config, remove it
	_ = uc.marketplace.Delete(ctx, workspace.ID, credentialID)

	return cred, nil
}
