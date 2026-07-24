package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	ctx, span := q.tracer.Start(ctx, "ListTemplates")
	defer span.End()
	return q.certs.ListTemplates(ctx, editionID)
}
