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

	// if this credential was backing the marketplace config, remove it
	config, err := uc.marketplace.Get(ctx, workspace.ID)
	if err == nil && config.CredentialID == credentialID {
		_ = uc.marketplace.Delete(ctx, workspace.ID)
	}

	return cred, nil
}
