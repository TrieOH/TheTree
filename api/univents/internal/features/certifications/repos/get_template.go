package repos

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetTemplateByID")
	defer span.End()
	tpl, err := database.Queries(ctx, repo.q).GetCertificationTemplateByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertificationTemplate(tpl)), nil
}
