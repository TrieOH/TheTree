package commands

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/internal/providers"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Checkout(ctx context.Context, payload models.CreateIntentInput) (*models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "Checkout")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	seller, err := c.sellers.GetByID(ctx, payload.SellerID)
	if err != nil {
		return nil, err
	}

	wallet, err := c.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return nil, err
	}
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleAdmin)
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
	provider, ok := providers.PayssageProviders.Payments[providerEnum]
	if !ok {
		return nil, fun.ErrBadRequest("provider not implemented")
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

	err = provider.Checkout(ctx, intent, payload.CheckoutData)
	if err != nil {
		return nil, err
	}

	created, err := c.intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
