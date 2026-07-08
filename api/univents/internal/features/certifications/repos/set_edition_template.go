package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) SetEditionTemplate(ctx context.Context, editionID uuid.UUID, templateID *uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "SetEditionTemplate")
	defer span.End()
	err := database.Queries(ctx, repo.q).SetEditionCertificationTemplate(ctx, sqlc.SetEditionCertificationTemplateParams{
		CertificationTemplateID: templateID,
		ID:                      editionID,
	})
	return repo.dbe(err)
}
