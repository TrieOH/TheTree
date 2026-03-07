package commands

import (
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) RevokeAPIKey(ctx context.Context, workspaceName string, keyID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RevokeAPIKey")
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
		Object("api_keys").
		Action("revoke").
		Scope(workspace.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("api key").SetMessage("insufficient permissions")
	}

	if _, err := uc.apiKeys.Revoke(ctx, keyID, workspace.ID); err != nil {
		return err
	}

	return nil
}
