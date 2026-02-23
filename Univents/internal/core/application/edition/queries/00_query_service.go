package queries

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	events   domain.EventsRepository
	editions domain.EditionsRepository
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	events domain.EventsRepository,
	editions domain.EditionsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		events:   events,
		editions: editions,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
