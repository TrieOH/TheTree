// Package intents is the intents feature: its repo, handlers and operations.
// Build wires it.
package intents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lib/telemetry"
	"log/slog"
	"payssage/errs"
	"payssage/internal/authz"
	"payssage/internal/providers"
	"payssage/models"
	"payssage/ports"
	idx "sdk/identityx"
	"uuid"
)

type Deps struct {
	Intents    ports.IntentRepo
	Wallets    ports.WalletRepo
	Orgs       ports.OrganizationRepo
	Collectors ports.CollectorRepo
	Sellers    ports.SellerRepo
	Authz      *authz.Service
	Registry   providers.PaymentProviders
	TestMode   bool // allows the /testmode routes
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

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
		telemetry.Log().ErrorContext(ctx, "error cancelling payment",
			slog.String("provider", intent.Provider),
			slog.String("intent_id", intent.ID.String()),
			telemetry.Err(err),
		)
		return nil, err
	}

	updated, err := o.Intents.Update(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

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

	id := uuid.NewV7()
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

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	intent, err := o.Intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, intent.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return intent, nil
}

func (o *Operations) HardCreate(ctx context.Context, req models.HardCreateIntentInput) (*models.Intent, error) {
	if !o.TestMode {
		return nil, errs.ErrTestModeOnly
	}
	_, err := o.Authz.Wallet(ctx, req.WalletID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

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

func (o *Operations) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Intent, error) {
	org, err := o.Orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Intents.ListByOrg(ctx, org.ID)
}

func (o *Operations) ListByProfile(ctx context.Context) ([]models.Intent, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.Intents.ListByOwner(ctx, ident.Sub.ID)
}

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error) {
	_, err := o.Authz.Wallet(ctx, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Intents.ListByWallet(ctx, walletID)
}

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
		telemetry.Log().ErrorContext(ctx, "error refunding payment",
			slog.String("provider", intent.Provider),
			slog.String("intent_id", intent.ID.String()),
			telemetry.Err(err),
		)
		return nil, err
	}

	updated, err := o.Intents.Update(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// TestmodeRefund simulates a provider refund (only when TestMode is on, and
// only for the wallet's admins): flips a succeeded intent to refunded and
// stamps a fake refund marker in provider_data — no provider is contacted.
// Mirrors the real refund's postconditions so downstream systems observe
// the same state the payment.refunded webhook would produce.
func (o *Operations) TestmodeRefund(ctx context.Context, intentID uuid.UUID) (*models.Intent, error) {
	if !o.TestMode {
		return nil, errs.ErrTestModeOnly
	}
	intent, err := o.Intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, intent.WalletID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}
	if intent.Status == models.IntentStatusRefunded {
		return intent, nil // idempotent
	}
	if intent.Status != models.IntentStatusSucceeded {
		return nil, &errs.IntentStatusError{Action: "refunded", Status: string(intent.Status)}
	}

	var pd models.MercadoPagoIntentData
	_ = json.Unmarshal(intent.ProviderData, &pd)
	fakeID := "testmode-refund-" + intent.ID.String()[:8]
	pd.RefundID = &fakeID
	status := "approved"
	pd.RefundStatus = &status
	pd.RefundAmountCents = &intent.AmountCents
	raw, err := json.Marshal(&pd)
	if err != nil {
		return nil, fmt.Errorf("marshal provider data: %w", err)
	}
	intent.ProviderData = raw
	intent.Status = models.IntentStatusRefunded

	return o.Intents.Update(ctx, *intent)
}
