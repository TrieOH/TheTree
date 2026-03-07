package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"

	"github.com/google/uuid"
)

func (uc *QueryService) GetByID(ctx context.Context, id uuid.UUID) (intent *domain.Intent, err error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.GetByID")
	defer span.End()

	ws, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err = uc.intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if intent.WorkspaceID != ws.ID {
		return nil, errx.Forbidden("intent").SetMessage("intent does not belong to this workspace")
	}

	return intent, nil
}
