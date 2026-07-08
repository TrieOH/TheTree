package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) ListByTarget(ctx context.Context, targetType string, targetID uuid.UUID) ([]contracts.Certification, error) {
	ctx, span := database.Span(ctx, repo.tracer, "ListByTarget")
	defer span.End()
	certs, err := database.Queries(ctx, repo.q).ListTargetCertifications(ctx, sqlc.ListTargetCertificationsParams{
		TargetType: targetType,
		TargetID:   targetID,
	})
	return xslices.MapSlice(certs, mapCertification), repo.dbe(err)
}
