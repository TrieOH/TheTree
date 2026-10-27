package service

import (
	"context"
	"errors"
	"payssage/errs"
	"payssage/internal/providers"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) Checkout(ctx context.Context, payload models.CreateIntentInput) (*models.Intent, error) {
	seller, err := o.Sellers.GetByID(ctx, payload.SellerID)
	if err != nil {
		return nil, err
	}

	wallet, err := o.Wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

	if seller.WalletID != wallet.ID {
		return nil, errs.ErrSellerNotInWallet
	}

	var collectorID *uuid.UUID
	if wallet.CollectorID != nil {
		collector, err := o.Collectors.GetByID(ctx, *wallet.CollectorID)
		if err != nil {
			return nil, err
		}
		if collector.Provider != seller.Provider {
			return nil, errors.New("seller and collector providers do not match")
		}
		collectorID = &collector.ID
	}

	providerEnum, err := providers.FromString(seller.Provider)
	if err != nil {
		return nil, err
	}
	provider, ok := o.Registry.Payments[providerEnum]
	if !ok {
		return nil, errs.ErrProviderNotImplemented
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	intent := &models.Intent{
		ID:            id,
		WalletID:      wallet.ID,
		SellerID:      seller.ID,
		CollectorID:   collectorID,
		AmountCents:   payload.AmountCents,
		Currency:      payload.Currency,
		Sandbox:       wallet.Sandbox,
		Provider:      seller.Provider,
		Status:        models.IntentStatusPending,
		Metadata:      payload.Metadata,
		ExternalID:    payload.ExternalID,
		ExternalGroup: payload.ExternalGroup,
	}

	err = provider.Checkout(ctx, intent, payload.CheckoutData)
	if err != nil {
		return nil, err
	}

	created, err := o.Intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
