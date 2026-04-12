package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *QueryService) ListWebhookEvents(ctx context.Context, workspaceName string) ([]domain.WebhookEventOriginal, error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.ListWebhookEvents")
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
		authz.Permission("view_webhook_events"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.events.ListByWorkspace(ctx, workspace.ID)
}
