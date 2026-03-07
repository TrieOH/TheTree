package commands

import (
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) DeleteWebhookEndpoint(ctx context.Context, workspaceName string, endpointID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DeleteWebhookEndpoint")
	defer span.End()

	ga := uc.gaClient

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("webhooks").
		Action("delete").
		Scope(workspace.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("webhooks").SetMessage("insufficient permissions")
	}

	return uc.endpoints.Delete(ctx, endpointID, workspace.ID)
}
