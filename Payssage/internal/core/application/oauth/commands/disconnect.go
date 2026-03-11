package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
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

	// Remove marketplace config for this credential if one exists
	if delErr := uc.marketplace.Delete(ctx, workspace.ID, credentialID); delErr != nil && !errx.IsKind(delErr, "not_found") {
		return nil, delErr
	}

	return uc.credentials.Revoke(ctx, credentialID, workspace.ID)
}
