package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListEmissionErrors(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListEmissionErrors")
	defer span.End()
	return q.certs.ListEmissionErrorsByEdition(ctx, editionID)
}
