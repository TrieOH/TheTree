package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListCertsByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListCertsByUser")
	defer span.End()
	return o.certs.ListByUser(ctx, userID)
}
