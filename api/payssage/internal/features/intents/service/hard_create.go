package service

import (
	"context"
	"payssage/models"
	"uuid"
)

func (o *Operations) HardCreate(ctx context.Context, req models.HardCreateIntentInput) (*models.Intent, error) {
	id := uuid.NewV7()

	intent := &models.Intent{
		ID:            id,
		WalletID:      req.WalletID,
		SellerID:      req.SellerID,
		CollectorID:   req.CollectorID,
		AmountCents:   req.AmountCents,
		Currency:      req.Currency,
		Sandbox:       req.Sandbox,
		Provider:      req.Provider,
		Status:        req.Status,
		ProviderData:  req.ProviderData,
		Metadata:      req.Metadata,
		ExternalID:    req.ExternalID,
		ExternalGroup: req.ExternalGroup,
	}

	return o.Intents.Create(ctx, *intent)
}
