package mercado_pago_provider

import (
	"lib/database"
	"lib/errx"
	"payssage/internal/config"
	"payssage/ports"

	mpconfig "github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/oauth"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var _ ports.PaymentAbstractionLayer = (*Provider)(nil)

type Provider struct {
	cfg         config.MercadoPagoConfig
	intents     ports.IntentRepo
	collectors  ports.CollectorRepo
	sellers     ports.SellerRepo
	wallets     ports.WalletRepo
	oauthClient oauth.Client
	logger      *zap.Logger
	tracer      trace.Tracer
	tx          database.TxRunner
}

func NewProvider(
	cfg config.MercadoPagoConfig,
	intents ports.IntentRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
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
		logger:      logger,
		tracer:      tracer,
		tx:          tx,
	}
}
