package service

import (
	"context"
	"lib/telemetry"
	"payssage/errs"
	"payssage/internal/providers"
	"payssage/models"

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

func (o *Operations) Cancel(ctx context.Context, intentID uuid.UUID) (*models.Intent, error) {
	intent, err := o.Intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, intent.WalletID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

	if !cancellableStatuses[intent.Status] {
		return nil, &errs.IntentStatusError{Action: "cancelled", Status: string(intent.Status)}
	}

	providerEnum, err := providers.FromString(intent.Provider)
	if err != nil {
		return nil, err
	}
	provider, ok := o.Registry.Payments[providerEnum]
	if !ok {
		return nil, errs.ErrProviderNotImplemented
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

	updated, err := o.Intents.Update(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
