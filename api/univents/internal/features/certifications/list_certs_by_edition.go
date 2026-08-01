package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListCertsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListCertsByEdition")
	defer span.End()
	return o.certs.ListByEdition(ctx, editionID)
}
