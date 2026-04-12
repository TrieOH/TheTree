package commands

import (
	"TriePayments/internal/shared/authz"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) DeleteWebhookEndpoint(ctx context.Context, workspaceName string, endpointID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DeleteWebhookEndpoint")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("delete_webhooks"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return err
	}

	return uc.endpoints.Delete(ctx, endpointID, workspace.ID)
}
