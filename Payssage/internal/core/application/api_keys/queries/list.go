package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *QueryService) List(ctx context.Context, workspaceName string) (ak []domain.APIKey, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspace *domain.Workspace
	workspace, err = uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	var keys []domain.APIKey
	keys, err = uc.apiKeys.ListByWorkspace(ctx, workspace.ID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
