package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"

	"github.com/google/uuid"
)

func (uc *QueryService) ListWebhookDeliveries(ctx context.Context, workspaceName string, endpointID uuid.UUID) ([]domain.WebhookDelivery, error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.ListWebhookDeliveries")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_webhook_deliveries"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.deliveries.ListByEndpoint(ctx, endpointID)
}
