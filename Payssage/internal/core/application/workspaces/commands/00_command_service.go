package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	workspaces domain.WorkspaceRepo
	gaClient   *goauth.Client
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func New(
	workspaces domain.WorkspaceRepo,
	gaClient *goauth.Client,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		workspaces: workspaces,
		gaClient:   gaClient,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}
