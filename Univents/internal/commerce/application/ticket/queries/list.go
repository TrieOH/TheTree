package queries

import (
	"context"
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
)

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []domain.Ticket, err error) { // FIXME Pagination
	return uc.tickets.List(ctx, editionID)
}
