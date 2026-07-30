package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListCertsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListCertsByEdition")
	defer span.End()
	return q.certs.ListByEdition(ctx, editionID)
}
