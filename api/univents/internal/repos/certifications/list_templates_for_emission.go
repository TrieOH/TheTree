package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListTemplatesForEmission(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListTemplatesForEmission")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListCertificationTemplatesByEditionForEmission(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertTemplate), nil
}
