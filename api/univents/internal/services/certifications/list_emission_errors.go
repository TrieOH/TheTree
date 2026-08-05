package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListEmissionErrors(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListEmissionErrors")
	defer span.End()
	return o.certs.ListEmissionErrorsByEdition(ctx, editionID)
}
