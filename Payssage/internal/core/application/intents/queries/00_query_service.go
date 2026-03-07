package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	intents    domain.IntentRepository
	workspaces domain.WorkspaceRepo
	gaClient   *goauth.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func New(
	intents domain.IntentRepository,
	workspaces domain.WorkspaceRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		intents:    intents,
		workspaces: workspaces,
		gaClient:   gaClient,
		tx:         tx,
		tracer:     tracer,
	}
}
