package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type CreateIntentInput struct {
	Amount             int64
	Currency           string
	Provider           string
	Metadata           json.RawMessage
	PaymentMethodID    string
	Installments       int
	CardToken          string
	PaymentMethodType  string
	SellerCredentialID uuid.UUID
	PayerEmail         string
}

func (uc *CommandService) InitiateCheckout(ctx context.Context, in CreateIntentInput) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.InitiateCheckout")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	marketplaceConfig, err := uc.marketplace.GetByProvider(ctx, workspace.ID, in.Provider)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			return nil, errx.Invalid("Provider").SetMessage(
				fmt.Sprintf("Provider '%s' is not configured for this workspace", in.Provider),
			)
		}
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

	var intent *domain.Intent

	switch p := provider.(type) {
	case domain.MercadoPagoProvider:
		var validationErrors []string
		if credential.Credentials.AccessToken == "" {
			validationErrors = append(validationErrors, "missing access token")
		}
		if in.PayerEmail == "" {
			validationErrors = append(validationErrors, "missing payer email")
		}
		if in.Amount <= 0 {
			validationErrors = append(validationErrors, "invalid amount")
		}
		if in.PaymentMethodType == "" {
			validationErrors = append(validationErrors, "missing payment method type")
		}

		if in.PaymentMethodType != "pix" {
			if in.Installments == 0 {
				validationErrors = append(validationErrors, "missing installments")
			}
			if in.PaymentMethodID == "" {
				validationErrors = append(validationErrors, "missing payment method id")
			}
			if in.CardToken == "" {
				validationErrors = append(validationErrors, "missing card token")
			}
		}

		if len(validationErrors) > 0 {
			return nil, errors.New("validation failed:\n" + strings.Join(validationErrors, "\n"))
		}

		baseRequest := &domain.InitiateCheckoutRequest{
			WorkspaceID:         workspace.ID,
			SellerCredentialID:  in.SellerCredentialID,
			Amount:              in.Amount,
			Currency:            strings.ToUpper(in.Currency),
			Provider:            in.Provider,
			Metadata:            in.Metadata,
			Payer:               domain.Payer{Email: in.PayerEmail},
			MPSellerToken:       credential.Credentials.AccessToken,
			MPMarketplaceFeeBPS: marketplaceConfig.FeeBps,
		}

		if in.PaymentMethodType == "pix" {
			intent, err = p.InitiatePixCheckout(ctx, baseRequest)
		} else {
			baseRequest.Installments = in.Installments
			baseRequest.MPPaymentMethodID = in.PaymentMethodID
			baseRequest.MPPaymentMethodType = in.PaymentMethodType
			baseRequest.MPPayerToken = in.CardToken
			intent, err = p.InitiateCheckout(ctx, baseRequest)
		}

		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown provider type: %T", p)
	}

	created, err := uc.intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
