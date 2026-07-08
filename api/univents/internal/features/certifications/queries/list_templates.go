package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]contracts.CertificationTemplate, error) {
	ctx, span := q.tracer.Start(ctx, "ListTemplates")
	defer span.End()
	return q.certs.ListTemplates(ctx, editionID)
}
