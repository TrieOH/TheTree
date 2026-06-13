package intents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"payssage/internal/platform/database"
	"payssage/internal/platform/telemetry"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	intents          ports.IntentRepository
	workspaces       ports.WorkspaceRepo
	credentials      ports.ProviderCredentialRepo
	marketplace      ports.MarketplaceConfigRepo
	webhooks         ports.WebhookDispatcher
	oauthProvider    map[string]ports.OAuthProvider
	paymentProviders map[string]ports.PaymentAbstractionLayer
	tx               database.TxRunner
	tracer           trace.Tracer
}

func NewCommandService(
	intents ports.IntentRepository,
	workspaces ports.WorkspaceRepo,
	credentials ports.ProviderCredentialRepo,
	marketplace ports.MarketplaceConfigRepo,
	webhooks ports.WebhookDispatcher,
	oauthProvider map[string]ports.OAuthProvider,
	paymentProviders map[string]ports.PaymentAbstractionLayer,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		intents:          intents,
		workspaces:       workspaces,
		credentials:      credentials,
		marketplace:      marketplace,
		webhooks:         webhooks,
		oauthProvider:    oauthProvider,
		paymentProviders: paymentProviders,
		tx:               tx,
		tracer:           tracer,
	}
}

func (uc *CommandService) CancelIntent(ctx context.Context, intentID uuid.UUID) (*contracts.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CancelIntent")
	defer span.End()

	ws, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := uc.intents.GetByID(ctx, intentID)
	if err != nil {
		return nil, err
	}

	if intent.WorkspaceID != ws.ID {
		return nil, errx.Forbidden("intent").SetMessage("intent does not belong to this workspace")
	}

	cancelled, err := uc.intents.Cancel(ctx, intentID)
	if err != nil {
		return nil, err
	}

	return cancelled, nil
}

type CancelPixInput struct {
	Provider           string
	IntentID           uuid.UUID
	SellerCredentialID uuid.UUID
}

func (uc *CommandService) CancelPix(ctx context.Context, in CancelPixInput) (*contracts.Intent, error) {
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
	case ports.MercadoPagoProvider:
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

func (uc *CommandService) Charge(ctx context.Context, intentID uuid.UUID, sellerCredentialID uuid.UUID) (*contracts.Intent, error) {
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

	if intent.Status != contracts.IntentStatusPending {
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

	chargedIntent, err := provider.Charge(ctx, &ports.ChargeRequest{
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

type CreateIntentInput struct {
	Amount               int64
	Currency             string
	Provider             string
	Metadata             json.RawMessage
	PaymentMethodID      string
	Installments         int
	CardToken            string
	PaymentMethodType    string
	SellerCredentialID   uuid.UUID
	PayerEmail           string
	IdentificationNumber string
	IdentificationType   string
}

func (uc *CommandService) InitiateCheckout(ctx context.Context, in CreateIntentInput) (*contracts.Intent, error) {
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

	var intent *contracts.Intent

	switch p := provider.(type) {
	case ports.MercadoPagoProvider:
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
		if in.PaymentMethodID == "" {
			validationErrors = append(validationErrors, "missing payment method type")
		}
		if in.PaymentMethodType == "" {
			validationErrors = append(validationErrors, "missing payment method id")
		}

		if in.PaymentMethodID != "pix" {
			if in.Installments == 0 {
				validationErrors = append(validationErrors, "missing installments")
			}

			if in.CardToken == "" {
				validationErrors = append(validationErrors, "missing card token")
			}
		}

		if len(validationErrors) > 0 {
			return nil, errors.New("validation failed:\n" + strings.Join(validationErrors, "\n"))
		}

		baseRequest := &ports.InitiateCheckoutRequest{
			WorkspaceID:          workspace.ID,
			SellerCredentialID:   in.SellerCredentialID,
			Amount:               in.Amount,
			Currency:             strings.ToUpper(in.Currency),
			Provider:             in.Provider,
			Metadata:             in.Metadata,
			IdentificationType:   in.IdentificationType,
			IdentificationNumber: in.IdentificationNumber,
			Payer:                ports.Payer{Email: in.PayerEmail},
			MPSellerToken:        credential.Credentials.AccessToken,
			MPMarketplaceFeeBPS:  marketplaceConfig.FeeBps,
		}

		if in.PaymentMethodID == "pix" {
			intent, err = p.InitiatePixCheckout(ctx, baseRequest)
		} else {
			baseRequest.Installments = in.Installments
			baseRequest.MPPaymentMethodID = in.PaymentMethodID
			baseRequest.MPPaymentMethodType = in.PaymentMethodType
			baseRequest.MPCardToken = in.CardToken
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
