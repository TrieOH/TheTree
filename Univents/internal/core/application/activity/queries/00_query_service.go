package queries

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	activities domain.ActivitiesRepository
	editions   domain.EditionsRepository
	gaClient   *goauth.Client
	tracer     trace.Tracer
	tx         database.TxRunner
}

func New(
	activities domain.ActivitiesRepository,
	editions domain.EditionsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		activities: activities,
		editions:   editions,
		gaClient:   gaClient,
		tracer:     tracer,
		tx:         tx,
	}
}
