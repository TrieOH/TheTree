package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsQueries.ListCertTemplateLinks")
	defer span.End()
	return o.Certs.ListCertTemplateLinks(ctx, templateID)
}
