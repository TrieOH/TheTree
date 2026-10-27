package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error) {
	results, err := database.Queries(ctx, repo.q).ListCertificationTemplateProgramsByTemplate(ctx, templateID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertTemplateProgram), nil
}
