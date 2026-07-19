package commands

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (c *Commands) HardCreate(ctx context.Context, req models.HardCreateIntentRequest) (*models.Intent, error) {
	ctx, span := c.tracer.Start(ctx, "HardCreate")
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

	return c.intents.Create(ctx, *intent)
}
