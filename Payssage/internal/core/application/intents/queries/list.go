package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
)

func (uc *QueryService) List(ctx context.Context) (intents []domain.Intent, err error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.List")
	defer span.End()

	// try workspace from API key first
	ws, err := authz.RequireWorkspace(ctx)
	if err == nil {
		return uc.intents.ListIntentsByWorkspace(ctx, ws.ID)
	}

	// fall back to user session — list all workspaces then all intents
	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspaces, err := uc.workspaces.List(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	for _, w := range workspaces {
		ws_intents, err := uc.intents.ListIntentsByWorkspace(ctx, w.ID)
		if err != nil {
			return nil, err
		}
		intents = append(intents, ws_intents...)
	}

	return intents, nil
}

func (uc *QueryService) ListByWorkspace(ctx context.Context, wsName string) (intents []domain.Intent, err error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.List")
	defer span.End()

	// try workspace from API key first
	ws, err := authz.RequireWorkspace(ctx)
	if err == nil {
		return uc.intents.ListIntentsByWorkspace(ctx, ws.ID)
	}

	// fall back to user session — list all workspaces then all intents
	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, wsName, sub.ID)
	if err != nil {
		return nil, err
	}

	intents, err = uc.intents.ListIntentsByWorkspace(ctx, workspace.ID)
	if err != nil {
		return nil, err
	}

	return intents, nil
}
