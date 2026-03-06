package queries

import (
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	tickets  domain.TicketsRepository
	editions domain2.EditionsRepository
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	tickets domain.TicketsRepository,
	editions domain2.EditionsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		tickets:  tickets,
		editions: editions,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
