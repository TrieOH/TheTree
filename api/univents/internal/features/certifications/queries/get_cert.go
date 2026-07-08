package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Certification, error) {
	ctx, span := q.tracer.Start(ctx, "GetByID")
	defer span.End()
	return q.certs.GetByID(ctx, id)
}
