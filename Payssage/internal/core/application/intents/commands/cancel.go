package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
)

func (uc *CommandService) CancelIntent(ctx context.Context, intentID uuid.UUID) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CancelIntent")
	defer span.End()

	ws, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := uc.intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}

	if intent.WorkspaceID != ws.ID {
		return nil, errx.Forbidden("intent").SetMessage("intent does not belong to this workspace")
	}

	cancelled, err := uc.intents.Cancel(ctx, intentID)
	if err != nil {
		return nil, err
	}

	return cancelled, nil
}
