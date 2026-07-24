package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	ctx, span := q.tracer.Start(ctx, "GetTemplateByID")
	defer span.End()
	return q.certs.GetTemplateByID(ctx, id)
}
