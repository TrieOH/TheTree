package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	intents          domain.IntentRepository
	workspaces       domain.WorkspaceRepo
	credentials      domain.ProviderCredentialRepo
	marketplace      domain.MarketplaceConfigRepo
	paymentProviders map[string]domain.PaymentProvider
	gaClient         *goauth.Client
	tx               database.TxRunner
	tracer           trace.Tracer
}

func New(
	intents domain.IntentRepository,
	workspaces domain.WorkspaceRepo,
	credentials domain.ProviderCredentialRepo,
	marketplace domain.MarketplaceConfigRepo,
	paymentProviders map[string]domain.PaymentProvider,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		intents:          intents,
		workspaces:       workspaces,
		credentials:      credentials,
		marketplace:      marketplace,
		paymentProviders: paymentProviders,
		gaClient:         gaClient,
		tx:               tx,
		tracer:           tracer,
	}
}
