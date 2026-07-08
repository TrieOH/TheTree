package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

func (q *Queries) GetTemplateByID(ctx context.Context, id, editionID uuid.UUID) (*contracts.CertificationTemplate, error) {
	ctx, span := q.tracer.Start(ctx, "GetTemplateByID")
	defer span.End()
	return q.certs.GetTemplateByID(ctx, id, editionID)
}
