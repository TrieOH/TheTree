package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	endpoints  domain.WebhookEndpointRepo
	deliveries domain.WebhookDeliveryRepo
	events     domain.WebhookEventRepo
	workspaces domain.WorkspaceRepo
	gaClient   *goauth.Client
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func New(
	endpoints domain.WebhookEndpointRepo,
	deliveries domain.WebhookDeliveryRepo,
	events domain.WebhookEventRepo,
	workspaces domain.WorkspaceRepo,
	gaClient *goauth.Client,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		endpoints:  endpoints,
		deliveries: deliveries,
		events:     events,
		workspaces: workspaces,
		gaClient:   gaClient,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}
