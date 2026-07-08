package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) ListByTarget(ctx context.Context, targetType string, targetID uuid.UUID) ([]contracts.Certification, error) {
	ctx, span := q.tracer.Start(ctx, "ListByTarget")
	defer span.End()
	return q.certs.ListByTarget(ctx, targetType, targetID)
}
