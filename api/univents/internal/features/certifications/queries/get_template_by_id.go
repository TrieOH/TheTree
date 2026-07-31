package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.GetTemplateByID")
	defer span.End()
	return q.certs.GetTemplateByID(ctx, id)
}
