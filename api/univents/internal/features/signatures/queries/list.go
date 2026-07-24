package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) List(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	ctx, span := q.tracer.Start(ctx, "List")
	defer span.End()

	return q.signatures.List(ctx, editionID)
}
