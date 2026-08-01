package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.GetTemplateByID")
	defer span.End()
	return o.certs.GetTemplateByID(ctx, id)
}
