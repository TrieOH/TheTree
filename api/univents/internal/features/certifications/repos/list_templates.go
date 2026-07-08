package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *repo) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]contracts.CertificationTemplate, error) {
	ctx, span := database.Span(ctx, repo.tracer, "ListTemplates")
	defer span.End()
	templates, err := database.Queries(ctx, repo.q).ListCertificationTemplates(ctx, editionID)
	return xslices.MapSlice(templates, mapCertificationTemplate), repo.dbe(err)
}
