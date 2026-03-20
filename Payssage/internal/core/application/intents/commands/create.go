package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/order"
	"github.com/spf13/viper"
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

func (uc *CommandService) CreateIntent(ctx context.Context, in CreateIntentInput) (*domain.Intent, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CreateIntent")
	defer span.End()

	workspace, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := domain.NewIntent(workspace.ID, in.Amount, in.Currency, in.Provider, in.Metadata)
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

	applicationFee := float64(intent.Amount) * float64(marketplaceConfig.FeeBps) / 10000 / 100.0
	amountInUnits := float64(intent.Amount) / 100.0

	var cfg *config.Config
	if viper.GetBool("TEST_MODE") {
		cfg, err = config.New(viper.GetString("MP_TEST_ACCESS_TOKEN"))
		if err != nil {
			return nil, err
		}
	} else {
		credential, err := uc.credentials.GetByID(ctx, in.SellerCredentialID)
		if err != nil {
			return nil, err
		}

		cfg, err = config.New(credential.Credentials.AccessToken)
		if err != nil {
			return nil, err
		}
	}

	client := order.NewClient(cfg)
	mpOrder, err := client.Create(ctx, order.Request{
		Type:              "online",
		TotalAmount:       fmt.Sprintf("%.2v", amountInUnits),
		ExternalReference: "tp_" + intent.ID.String(),
		ExpirationTime:    "",
		Currency:          in.Currency,
		MarketPlaceFee:    fmt.Sprintf("%.2v", applicationFee),
		Transactions: &order.TransactionRequest{
			Payments: []order.PaymentRequest{
				{
					Amount:         fmt.Sprintf("%.2v", amountInUnits),
					ExpirationTime: "",
					PaymentMethod: &order.PaymentMethodRequest{
						ID:                  in.PaymentMethodID,
						Type:                in.PaymentMethodType,
						Token:               in.CardToken,
						StatementDescriptor: "",
						Installments:        in.Installments,
					},
				},
			},
		},
		Payer: &order.PayerRequest{
			Email: in.PayerEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	intent.AddExtOrderID(mpOrder.ID)

	created, err := uc.intents.Create(ctx, *intent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
