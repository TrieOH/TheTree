package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.ListByEdition")
	defer span.End()
	return q.signatures.ListByEdition(ctx, editionID)
}
