package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) ListByUser(ctx context.Context, userID uuid.UUID) ([]contracts.Certification, error) {
	ctx, span := q.tracer.Start(ctx, "ListByUser")
	defer span.End()
	return q.certs.ListByUser(ctx, userID)
}
