package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) SetActivityTemplate(ctx context.Context, activityID uuid.UUID, templateID *uuid.UUID) error {
	ctx, span := database.Span(ctx, repo.tracer, "SetActivityTemplate")
	defer span.End()
	err := database.Queries(ctx, repo.q).SetActivityCertificationTemplate(ctx, sqlc.SetActivityCertificationTemplateParams{
		CertificationTemplateID: templateID,
		ID:                      activityID,
	})
	return repo.dbe(err)
}
