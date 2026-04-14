package tickets

import (
	"context"
	"univents/internal/platform/database"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	tickets  ports.TicketsRepository
	editions ports.EditionsRepository
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueryService(
	tickets ports.TicketsRepository,
	editions ports.EditionsRepository,
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

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []contracts.Ticket, err error) { // FIXME Pagination
	return uc.tickets.List(ctx, editionID)
}
