package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) List(ctx context.Context, editionID uuid.UUID) ([]contracts.Signature, error) {
	ctx, span := q.tracer.Start(ctx, "List")
	defer span.End()

	return q.signatures.List(ctx, editionID)
}
