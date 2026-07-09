package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*contracts.CertificationTemplate, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetTemplateByID")
	defer span.End()
	tpl, err := database.Queries(ctx, repo.q).GetCertificationTemplateByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertificationTemplate(tpl)), nil
}
