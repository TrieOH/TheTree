package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id, editionID uuid.UUID) (*contracts.Signature, error) {
	ctx, span := q.tracer.Start(ctx, "GetByID")
	defer span.End()

	return q.signatures.GetByID(ctx, id, editionID)
}
