package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) DisconnectCredential(ctx context.Context, credentialID uuid.UUID) (*domain.ProviderCredential, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DisconnectCredential")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	return uc.credentials.Revoke(ctx, credentialID, workspace.ID)
}
