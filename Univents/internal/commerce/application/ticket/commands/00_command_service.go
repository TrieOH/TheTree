package commands

import (
	"univents/internal/commerce/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	tickets  domain.TicketsRepository
	asynq    *asynq.Client
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	tickets domain.TicketsRepository,
	asynq *asynq.Client,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		tickets:  tickets,
		asynq:    asynq,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
