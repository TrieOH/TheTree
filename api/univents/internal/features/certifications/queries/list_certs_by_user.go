package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListCertsByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListCertsByUser")
	defer span.End()
	return q.certs.ListByUser(ctx, userID)
}
