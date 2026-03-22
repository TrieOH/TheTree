package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	endpoints   domain.WebhookEndpointRepo
	deliveries  domain.WebhookDeliveryRepo
	events      domain.WebhookEventRepo
	workspaces  domain.WorkspaceRepo
	intents     domain.IntentRepository
	credentials domain.ProviderCredentialRepo
	asynq       *asynq.Client
	gaClient    *goauth.Client
	tx          database.TxRunner
	tracer      trace.Tracer
}

func New(
	endpoints domain.WebhookEndpointRepo,
	deliveries domain.WebhookDeliveryRepo,
	events domain.WebhookEventRepo,
	workspaces domain.WorkspaceRepo,
	intents domain.IntentRepository,
	credentials domain.ProviderCredentialRepo,
	asynq *asynq.Client,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		endpoints:   endpoints,
		deliveries:  deliveries,
		events:      events,
		workspaces:  workspaces,
		intents:     intents,
		credentials: credentials,
		asynq:       asynq,
		gaClient:    gaClient,
		tx:          tx,
		tracer:      tracer,
	}
}
