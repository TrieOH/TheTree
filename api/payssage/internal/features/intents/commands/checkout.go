package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"payssage/internal/providers"
	"payssage/models"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Checkout(ctx context.Context, payload models.CreateIntentInput) (*models.Intent, error) {
	ctx, span := c.tracer.Start(ctx, "Checkout")
	defer span.End()

	seller, err := c.sellers.GetByID(ctx, payload.SellerID)
	if err != nil {
		return nil, err
	}

	wallet, err := c.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return nil, err
	}

	if seller.WalletID != wallet.ID {
		return nil, fun.ErrBadRequest("seller does not belong to wallet")
	}

	var collectorID *uuid.UUID
	if wallet.CollectorID != nil {
		collector, err := c.collectors.GetByID(ctx, *wallet.CollectorID)
		if err != nil {
			return nil, err
		}
		if collector.Provider != seller.Provider {
			return nil, fun.ErrInternal("seller and collector providers do not match")
		}
		collectorID = &collector.ID
	}

	providerEnum, err := providers.FromString(seller.Provider)
	if err != nil {
		return nil, err
	}

	payssageProvider, ok := providers.PayssageProviders.Payments[providerEnum]
	if !ok {
		return nil, fun.ErrBadRequest("invalid provider")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	intent := &models.Intent{
		ID:          id,
		WalletID:    wallet.ID,
		SellerID:    seller.ID,
		CollectorID: collectorID,
		AmountCents: payload.AmountCents,
		Currency:    payload.Currency,
		Sandbox:     wallet.Sandbox,
		Provider:    seller.Provider,
		Status:      models.IntentStatusPending,
		Metadata:    payload.Metadata,
	}

	var providerData json.RawMessage
	switch provider := payssageProvider.(type) {
	case ports.MercadoPagoProvider:
		providerData, err = provider.Checkout(ctx, intent, payload.CheckoutData)
		if err != nil {
			return nil, err
		}
		intent.ProviderData = providerData

	default:
		return nil, fmt.Errorf("unknown provider type: %T", provider)
	}

	created, err := c.intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
