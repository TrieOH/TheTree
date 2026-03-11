package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
)

type SetMarketplaceConfigRequest struct {
	WorkspaceName string
	CredentialID  uuid.UUID
	FeeBps        int
}

func (uc *CommandService) SetMarketplaceConfig(ctx context.Context, req SetMarketplaceConfigRequest) (*domain.MarketplaceConfig, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.SetMarketplaceConfig")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	// verify credential belongs to this workspace
	cred, err := uc.credentials.GetByID(ctx, req.CredentialID)
	if err != nil {
		return nil, err
	}
	if cred.WorkspaceID != workspace.ID {
		return nil, errx.Forbidden("credential").SetMessage("credential does not belong to this workspace")
	}

	existing, err := uc.marketplace.Get(ctx, workspace.ID, req.CredentialID)
	if err != nil && !errx.IsKind(err, "not_found") {
		return nil, err
	}

	if existing != nil {
		return uc.marketplace.Update(ctx, domain.MarketplaceConfig{
			WorkspaceID:  workspace.ID,
			CredentialID: req.CredentialID,
			FeeBps:       req.FeeBps,
		})
	}

	return uc.marketplace.Create(ctx, domain.MarketplaceConfig{
		WorkspaceID:  workspace.ID,
		CredentialID: req.CredentialID,
		FeeBps:       req.FeeBps,
	})
}
