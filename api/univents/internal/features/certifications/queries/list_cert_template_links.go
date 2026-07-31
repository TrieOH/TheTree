package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListCertTemplateLinks")
	defer span.End()
	return q.certs.ListCertTemplateLinks(ctx, templateID)
}
