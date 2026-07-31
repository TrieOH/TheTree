package mercado_pago_provider

import (
	"lib/errx"
	"payssage/internal/config"
	"payssage/ports"

	mpconfig "github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
	"resty.dev/v3"
)

var _ ports.PaymentAbstractionLayer = (*Provider)(nil)

type Provider struct {
	cfg         config.MercadoPagoConfig
	intents     ports.IntentRepo
	collectors  ports.CollectorRepo
	sellers     ports.SellerRepo
	wallets     ports.WalletRepo
	oauthClient oauth.Client
	httpClient  *resty.Client
}

func NewProvider(
	cfg config.MercadoPagoConfig,
	intents ports.IntentRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	httpClient *resty.Client,
) *Provider {
	mpCfg, err := mpconfig.New(cfg.MpAccessToken)
	if err != nil {
		errx.Exit(err, "error creating mercado pago provider")
	}
	return &Provider{
		cfg:         cfg,
		intents:     intents,
		collectors:  collectors,
		sellers:     sellers,
		wallets:     wallets,
		oauthClient: oauth.NewClient(mpCfg),
		httpClient:  httpClient,
	}
}
