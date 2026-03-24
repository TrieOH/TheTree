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

type CancelPixInput struct {
	Provider           string
	IntentID           uuid.UUID
	SellerCredentialID uuid.UUID
}

func (uc *CommandService) CancelPix(ctx context.Context, in CancelPixInput) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CancelPix")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	_, err = uc.marketplace.GetByProvider(ctx, workspace.ID, in.Provider)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			return nil, errx.Invalid("Provider").SetMessage(
				fmt.Sprintf("Provider '%s' is not configured for this workspace", in.Provider),
			)
		}
		return nil, err
	}

	intent, err := uc.intents.GetByID(ctx, in.IntentID)
	if err != nil {
		return nil, err
	}

	credential, err := uc.credentials.GetByID(ctx, in.SellerCredentialID)
	if err != nil {
		return nil, err
	}

	provider, ok := uc.paymentProviders[in.Provider]
	if !ok {
		return nil, errx.Invalid("No such payment provider")
	}

	switch p := provider.(type) {
	case domain.MercadoPagoProvider:
		err = p.CancelPixCode(ctx, intent.MercadoPagoData.TransactionID, credential.Credentials.AccessToken)
		if err != nil {
			telemetry.Log().Error("Error canceling PIX", zap.String("provider", in.Provider), zap.Error(err))
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown provider type: %T", p)
	}

	created, err := uc.intents.Cancel(ctx, intent.ID)
	if err != nil {
		return nil, err
	}

	return created, nil
}
