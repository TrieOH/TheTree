package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	ctx, span := q.tracer.Start(ctx, "GetByID")
	defer span.End()
	return q.certs.GetByID(ctx, id)
}
