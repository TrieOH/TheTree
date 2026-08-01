package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListTemplates")
	defer span.End()
	return o.certs.ListTemplates(ctx, editionID)
}
