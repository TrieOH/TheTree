package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
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

	allowed, err := uc.gaClient.Authz.Check().User(sub.ID).
		Object("webhooks").
		Action("read").
		Scope(workspace.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("webhooks").SetMessage("insufficient permissions")
	}

	return uc.deliveries.ListByEndpoint(ctx, endpointID)
}
