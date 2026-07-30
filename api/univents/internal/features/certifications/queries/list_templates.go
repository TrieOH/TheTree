package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListTemplates")
	defer span.End()
	return q.certs.ListTemplates(ctx, editionID)
}
