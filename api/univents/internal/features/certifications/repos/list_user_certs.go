package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *repo) ListByUser(ctx context.Context, userID uuid.UUID) ([]contracts.Certification, error) {
	ctx, span := database.Span(ctx, repo.tracer, "ListByUser")
	defer span.End()
	certs, err := database.Queries(ctx, repo.q).ListUserCertifications(ctx, userID)
	return xslices.MapSlice(certs, mapCertification), repo.dbe(err)
}
