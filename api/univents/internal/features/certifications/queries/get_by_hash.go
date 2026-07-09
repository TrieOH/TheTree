package queries

import (
	"context"
	"univents/contracts"
)

func (q *Queries) GetByHash(ctx context.Context, hash string) (*contracts.Certification, error) {
	ctx, span := q.tracer.Start(ctx, "GetByHash")
	defer span.End()
	return q.certs.GetByHash(ctx, hash)
}
