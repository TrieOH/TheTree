package repos

import (
	"context"
	"lib/crypto"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) Certify(ctx context.Context, input contracts.CertifyInput) (*contracts.Certification, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Certify")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	row, err := database.Queries(ctx, repo.q).Certify(ctx, sqlc.CertifyParams{
		ID:         id,
		UserID:     input.UserID,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
		Hash:       new(crypto.HashHMACSHA256(id.String())),
	})
	return new(mapCertification(row)), repo.dbe(err)
}
