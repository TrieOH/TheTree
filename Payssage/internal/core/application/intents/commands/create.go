package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"
	"fmt"
)

func (uc *CommandService) CreateIntent(ctx context.Context, amount int64, currency, provider string, metadata json.RawMessage) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CreateIntent")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	_, err = uc.credentials.GetByWorkspaceAndProvider(ctx, workspace.ID, provider)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			return nil, errx.Invalid("provider").SetMessage(
				fmt.Sprintf("provider '%s' is not configured for this workspace", provider),
			)
		}
		return nil, err
	}

	intent, err := domain.NewIntent(workspace.ID, amount, currency, provider, metadata)
	if err != nil {
		return nil, err
	}

	created, err := uc.intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
