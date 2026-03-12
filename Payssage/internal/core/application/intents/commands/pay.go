package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/telemetry"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

	if workspace.Sandbox {
		fakeProviderID := "sandbox_" + uuid.NewString()
		updated, err := uc.intents.Pay(ctx, intentID, fakeProviderID, domain.IntentStatusSucceeded)
		if err != nil {
			return nil, err
		}
		go func() {
			if err := uc.webhooks.Dispatch(
				context.Background(), "sandbox", updated.ID.String(), domain.EventPaymentSucceeded,
			); err != nil {
				log.Printf("[sandbox] webhook fan-out failed for intent=%s: %v", updated.ID, err)
			}
		}()
		return updated, nil
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

	marketplaceConfig, err := uc.marketplace.GetByProvider(ctx, workspace.ID, intent.Provider)
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
		telemetry.Log().Info("Charge Failed", zap.Error(err))
		_, _ = uc.intents.Fail(ctx, intentID)
		return nil, errx.Internal("payment").SetMessage("charge failed").SetCause(err)
	}

	updated, err := uc.intents.Pay(ctx, intentID, result.ProviderPaymentID, result.Status)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
