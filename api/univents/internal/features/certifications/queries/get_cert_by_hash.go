package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (q *Queries) GetCertByHash(ctx context.Context, hash string) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.GetCertByHash")
	defer span.End()
	return q.certs.GetByHash(ctx, hash)
}
