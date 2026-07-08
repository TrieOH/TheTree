package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Certification, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByID")
	defer span.End()
	cert, err := database.Queries(ctx, repo.q).GetCertificationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(cert)), nil
}
