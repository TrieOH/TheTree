package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListCertTemplateLinks")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListCertificationTemplateProgramsByTemplate(ctx, templateID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertTemplateProgram), nil
}
