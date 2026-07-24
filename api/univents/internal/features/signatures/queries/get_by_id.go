package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id, editionID uuid.UUID) (*models.Signature, error) {
	ctx, span := q.tracer.Start(ctx, "GetByID")
	defer span.End()

	return q.signatures.GetByID(ctx, id, editionID)
}
