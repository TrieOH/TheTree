package queries

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	checkpoints domain.CheckpointsRepository
	editions    domain.EditionsRepository
	gaClient    *goauth.Client
	tracer      trace.Tracer
	tx          database.TxRunner
}

func New(
	checkpoints domain.CheckpointsRepository,
	editions domain.EditionsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		checkpoints: checkpoints,
		editions:    editions,
		gaClient:    gaClient,
		tracer:      tracer,
		tx:          tx,
	}
}
