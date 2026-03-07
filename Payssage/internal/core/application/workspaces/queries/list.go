package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *QueryService) List(ctx context.Context) (ws []domain.Workspace, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.List")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspaces []domain.Workspace
	workspaces, err = uc.workspaces.List(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}
