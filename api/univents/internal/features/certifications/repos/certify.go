package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *repo) Certify(ctx context.Context, input contracts.CertifyInput) (*contracts.Certification, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Certify")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).Certify(ctx, sqlc.CertifyParams{
		UserID:     input.UserID,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(row)), nil
}
