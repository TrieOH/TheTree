package commands

import (
	"TriePayments/internal/core/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type PayIntentInput struct {
	CardToken          string
	PaymentMethodID    string
	Installments       int
	PayerEmail         string
	SellerCredentialID uuid.UUID
}

func (uc *CommandService) PayIntent(ctx context.Context, intentID uuid.UUID, input PayIntentInput) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.PayIntent")
	defer span.End()

	return &domain.Intent{}, errors.New("not implemented")
	/*
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
			updated, err := uc.intents.Confirm(ctx, intent.ID)
			if err != nil {
				return nil, err
			}
			go func() {
				if err := uc.webhooks.Dispatch(
					context.Background(), "sandbox", updated.ID.String(), domain.EventPaymentSucceeded, uuid.Nil,
				); err != nil {
					log.Printf("[sandbox] webhook fan-out failed for intent=%s: %v", updated.ID, err)
				}
			}()
			return updated, nil
		}

		provider, ok := uc.paymentProviders[intent.Provider]
		if !ok {
			return nil, errx.Invalid("Provider").SetMessage(
				fmt.Sprintf("payment Provider '%s' is not supported", intent.Provider),
			)
		}

		var token string
		var credential *domain.ProviderCredential
		if viper.GetBool("TEST_MODE") {
			token = viper.GetString("MP_TEST_ACCESS_TOKEN")
		} else {
			credential, err = uc.credentials.GetByID(ctx, input.SellerCredentialID)
			if err != nil {
				return nil, err
			}

			token = credential.Credentials.AccessToken
		}

		marketplaceConfig, err := uc.marketplace.GetByProvider(ctx, workspace.ID, intent.Provider)
		if err != nil {
			return nil, err
		}

		applicationFee := float64(intent.Amount) * float64(marketplaceConfig.FeeBps) / 10000 / 100.0
		amountInUnits := float64(intent.Amount) / 100.0

		result, err := provider.Charge(ctx, domain.ChargeRequestOriginal{
			SellerToken:     token,
			Installments:    input.Installments,
			CardToken:       input.CardToken,
			PayerEmail:      input.PayerEmail,
			Amount:          amountInUnits,
			PaymentMethodID: input.PaymentMethodID,
			ApplicationFee:  applicationFee,
			IntentID:        intent.ID,
		})
		if err != nil {
			// Mark intent as failed so it can't be retried with a stale token
			telemetry.Log().Info("Charge Failed", zap.Error(err))
			_, _ = uc.intents.Fail(ctx, intentID)
			return nil, errx.Internal("payment").SetMessage("charge failed").SetCause(err)
		}

		updated, err := uc.intents.Pay(ctx, intentID, result.OrderID, result.Status)
		if err != nil {
			return nil, err
		}

		return updated, nil*/
}
