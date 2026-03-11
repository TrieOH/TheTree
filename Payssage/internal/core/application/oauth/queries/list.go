package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *QueryService) ListMarketplaceConfigs(ctx context.Context, workspaceName string) ([]domain.MarketplaceConfig, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.ListMarketplaceConfigs")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	return uc.marketplaces.List(ctx, workspace.ID)
}
