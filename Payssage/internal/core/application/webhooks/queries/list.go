package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *QueryService) ListWebhookEndpoints(ctx context.Context, workspaceName string) ([]domain.WebhookEndpoint, error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.ListWebhookEndpoints")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	return uc.endpoints.ListByWorkspace(ctx, workspace.ID)
}
