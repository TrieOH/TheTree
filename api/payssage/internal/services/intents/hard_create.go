package intents

import (
	"context"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) HardCreate(ctx context.Context, req models.HardCreateIntentRequest) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "HardCreate")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	intent := &models.Intent{
		ID:           id,
		WalletID:     req.WalletID,
		SellerID:     req.SellerID,
		CollectorID:  req.CollectorID,
		AmountCents:  req.AmountCents,
		Currency:     req.Currency,
		Sandbox:      req.Sandbox,
		Provider:     req.Provider,
		Status:       req.Status,
		ProviderData: req.ProviderData,
		Metadata:     req.Metadata,
	}

	return o.intents.Create(ctx, *intent)
}
