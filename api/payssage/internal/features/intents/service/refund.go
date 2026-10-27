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

// Refund reverses an approved (succeeded) payment intent — full refund only
// in v1 (no amount). The provider records the refund in provider_data but
// the intent status stays `succeeded`; the payment.refunded webhook is what
// flips it (webhook-only confirmation, mirroring approval, D3).
func (o *Operations) Refund(ctx context.Context, intentID uuid.UUID) (*models.Intent, error) {
	intent, err := o.Intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, intent.WalletID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

	// Guard: only succeeded intents are refundable. An already-refunded
	// intent is an idempotent no-op (duplicate organizer clicks get the
	// current intent back, never a second refund).
	if intent.Status == models.IntentStatusRefunded {
		return intent, nil
	}
	if intent.Status != models.IntentStatusSucceeded {
		return nil, &errs.IntentStatusError{Action: "refunded", Status: string(intent.Status)}
	}

	providerEnum, err := providers.FromString(intent.Provider)
	if err != nil {
		return nil, err
	}
	provider, ok := o.Registry.Payments[providerEnum]
	if !ok {
		return nil, errs.ErrProviderNotImplemented
	}

	err = provider.Refund(ctx, intent)
	if err != nil {
		telemetry.Log().Error("error refunding payment",
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
