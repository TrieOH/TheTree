package commands

import (
	"univents/internal/commerce/domain"
	coreDomain "univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	editions coreDomain.EditionsRepository
	tickets  domain.TicketsRepository
	asynq    *asynq.Client
	gaClient *goauth.Client
	tracer   trace.Tracer
	az       *authzed.Client
	tx       database.TxRunner
}

func New(
	editions coreDomain.EditionsRepository,
	tickets domain.TicketsRepository,
	asynq *asynq.Client,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		editions: editions,
		tickets:  tickets,
		asynq:    asynq,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
