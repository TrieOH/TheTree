package certifications

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) GetCertByHash(ctx context.Context, hash string) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.GetCertByHash")
	defer span.End()
	return o.certs.GetByHash(ctx, hash)
}
