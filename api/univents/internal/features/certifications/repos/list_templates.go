package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListTemplates")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListCertificationTemplatesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertTemplate), nil
}
