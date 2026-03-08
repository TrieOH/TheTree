package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PayIntentInput struct {
	CardToken       string
	PaymentMethodID string
	Installments    int
	PayerEmail      string
}

func (uc *CommandService) PayIntent(ctx context.Context, intentID uuid.UUID, input PayIntentInput) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.PayIntent")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := uc.intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}

	if intent.WorkspaceID != workspace.ID {
		return nil, errx.NotFound("intent")
	}

	if intent.Status != domain.IntentStatusPending {
		return nil, errx.Invalid("intent").SetMessage("intent is not in a payable state")
	}

	provider, ok := uc.paymentProviders[intent.Provider]
	if !ok {
		return nil, errx.Invalid("provider").SetMessage(
			fmt.Sprintf("payment provider '%s' is not supported", intent.Provider),
		)
	}

	credential, err := uc.credentials.GetByWorkspaceAndProvider(ctx, workspace.ID, intent.Provider)
	if err != nil {
		return nil, err
	}

	marketplaceConfig, err := uc.marketplace.Get(ctx, workspace.ID)
	if err != nil {
		return nil, err
	}

	applicationFee := float64(intent.Amount) * float64(marketplaceConfig.FeeBps) / 10000 / 100.0

	result, err := provider.Charge(ctx, domain.ChargeRequest{
		Intent:          *intent,
		CardToken:       input.CardToken,
		PaymentMethodID: input.PaymentMethodID,
		Installments:    input.Installments,
		PayerEmail:      input.PayerEmail,
		ApplicationFee:  applicationFee,
		SellerToken:     credential.Credentials.AccessToken,
	})
	if err != nil {
		// Mark intent as failed so it can't be retried with a stale token
		_, _ = uc.intents.Fail(ctx, intentID)
		return nil, errx.Internal("payment").SetMessage("charge failed").SetCause(err)
	}

	updated, err := uc.intents.Pay(ctx, intentID, result.ProviderPaymentID, result.Status)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
