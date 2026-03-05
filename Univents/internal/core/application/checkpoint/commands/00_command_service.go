package commands

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	checkpoints domain.CheckpointsRepository
	editions    domain.EditionsRepository
	asynq       *asynq.Client
	gaClient    *goauth.Client
	tracer      trace.Tracer
	tx          database.TxRunner
}

func New(
	checkpoints domain.CheckpointsRepository,
	editions domain.EditionsRepository,
	asynq *asynq.Client,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		checkpoints: checkpoints,
		editions:    editions,
		asynq:       asynq,
		gaClient:    gaClient,
		tracer:      tracer,
		tx:          tx,
	}
}
