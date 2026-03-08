package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	intents     domain.IntentRepository
	workspaces  domain.WorkspaceRepo
	oauthStates domain.OAuthStateRepo
	credentials domain.ProviderCredentialRepo
	providers   map[string]domain.OAuthProvider
	gaClient    *goauth.Client
	tx          database.TxRunner
	tracer      trace.Tracer
}

func New(
	intents domain.IntentRepository,
	workspaces domain.WorkspaceRepo,
	oauthStates domain.OAuthStateRepo,
	credentials domain.ProviderCredentialRepo,
	providers map[string]domain.OAuthProvider,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		intents:     intents,
		workspaces:  workspaces,
		oauthStates: oauthStates,
		credentials: credentials,
		providers:   providers,
		gaClient:    gaClient,
		tx:          tx,
		tracer:      tracer,
	}
}
