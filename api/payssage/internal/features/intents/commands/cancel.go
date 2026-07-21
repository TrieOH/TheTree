package commands

import (
	"context"
	"lib/telemetry"
	"payssage/internal/providers"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// cancellableStatuses is the set of intent statuses that can still be
// canceled. Anything outside this set has already reached a terminal
// or provider-settled state and cannot be canceled.
var cancellableStatuses = map[models.IntentStatus]bool{
	models.IntentStatusPending:    true,
	models.IntentStatusProcessing: true,
}

func (c *Commands) Cancel(ctx context.Context, intentID uuid.UUID) (*models.Intent, error) {
	ctx, span := c.tracer.Start(ctx, "CancelIntent")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := c.intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}

	if err := c.checkAdminAccess(ctx, intent.WalletID, ident.Sub.ID); err != nil {
		return nil, err
	}

	if !cancellableStatuses[intent.Status] {
		return nil, fun.Errf("intent cannot be cancelled from status %q", intent.Status).BadRequest()
	}

	providerEnum, err := providers.FromString(intent.Provider)
	if err != nil {
		return nil, err
	}
	provider, ok := providers.PayssageProviders.Payments[providerEnum]
	if !ok {
		return nil, fun.Err("provider not implemented").BadRequest()
	}

	err = provider.CancelPendingPayment(ctx, intent)
	if err != nil {
		telemetry.Log().Error("error cancelling payment",
			zap.String("provider", intent.Provider),
			zap.String("intent_id", intent.ID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	updated, err := c.intents.Update(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
