package queries

import (
	"context"
	"univents/models"
)

func (q *Queries) GetByHash(ctx context.Context, hash string) (*models.Certification, error) {
	ctx, span := q.tracer.Start(ctx, "GetByHash")
	defer span.End()
	return q.certs.GetByHash(ctx, hash)
}
