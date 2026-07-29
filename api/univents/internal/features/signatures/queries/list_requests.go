package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.ListRequestsByEdition")
	defer span.End()
	return q.requests.ListRequestsByEdition(ctx, editionID)
}
