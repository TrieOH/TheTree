package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
	"encoding/json"
)

func (uc *CommandService) CreateIntent(ctx context.Context, amount int64, currency, provider string, metadata json.RawMessage) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CreateIntent")
	defer span.End()

	ws, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := domain.NewIntent(ws.ID, amount, currency, provider, metadata)
	if err != nil {
		return nil, err
	}

	created, err := uc.intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
