package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	apiKeys    domain.ApiKeysRepo
	workspaces domain.WorkspaceRepo
	gaClient   *goauth.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func New(
	apiKeys domain.ApiKeysRepo,
	workspaces domain.WorkspaceRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		apiKeys:    apiKeys,
		workspaces: workspaces,
		gaClient:   gaClient,
		tx:         tx,
		tracer:     tracer,
	}
}
