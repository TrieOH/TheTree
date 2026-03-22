package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/telemetry"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uc *CommandService) Charge(ctx context.Context, intentID uuid.UUID, sellerCredentialID uuid.UUID) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Charge")
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
		return nil, errx.Forbidden("intent").SetMessage("intent doesn't belong to this workspace")
	}

	if intent.Status != domain.IntentStatusPending {
		return nil, errx.Invalid("intent").SetMessage("intent is not in a payable state")
	}

	credential, err := uc.credentials.GetByID(ctx, sellerCredentialID)
	if err != nil {
		return nil, err
	}

	if credential.WorkspaceID != workspace.ID {
		return nil, errx.Forbidden("credential").SetMessage("credential doesn't belong to this workspace")
	}

	provider, ok := uc.paymentProviders[intent.Provider]
	if !ok {
		return nil, errx.Invalid("provider").SetMessage(fmt.Sprintf("payment provider '%s' is not supported", intent.Provider))
	}

	chargedIntent, err := provider.Charge(ctx, &domain.ChargeRequest{
		Intent:        *intent,
		MPSellerToken: credential.Credentials.AccessToken,
	})
	if err != nil {
		// Mark intent as failed so it can't be retried
		telemetry.Log().Info("Charge Failed", zap.Error(err))
		_, _ = uc.intents.Fail(ctx, intentID)
		return nil, errx.Internal("payment").SetMessage("charge failed").SetCause(err)
	}

	if _, err = uc.intents.UpdateProviderData(ctx, *chargedIntent); err != nil {
		return nil, err
	}

	return chargedIntent, nil
}
