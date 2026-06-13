package tickets

import (
	"context"

	"lib/database"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	tickets  ports.TicketsRepository
	editions ports.EditionsRepository
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueryService(
	tickets ports.TicketsRepository,
	editions ports.EditionsRepository,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		tickets:  tickets,
		editions: editions,
		tracer:   tracer,
		tx:       tx,
	}
}

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []contracts.Ticket, err error) { // FIXME Pagination
	return uc.tickets.List(ctx, editionID)
}
