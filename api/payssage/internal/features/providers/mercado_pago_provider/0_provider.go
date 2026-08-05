package mercado_pago_provider

import (
	"errors"
	"lib/errx"
	"payssage/internal/config"
	"payssage/ports"

	"resty.dev/v3"
)

var _ ports.PaymentAbstractionLayer = (*Provider)(nil)

type Provider struct {
	cfg        config.MercadoPagoConfig
	intents    ports.IntentRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
	wallets    ports.WalletRepo
	httpClient *resty.Client
}

func NewProvider(
	cfg config.MercadoPagoConfig,
	intents ports.IntentRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	httpClient *resty.Client,
) *Provider {
	if cfg.MpAccessToken == "" {
		errx.Exit(errors.New("missing mercado pago access token"), "error creating mercado pago provider")
	}
	return &Provider{
		cfg:        cfg,
		intents:    intents,
		collectors: collectors,
		sellers:    sellers,
		wallets:    wallets,
		httpClient: httpClient,
	}
}
