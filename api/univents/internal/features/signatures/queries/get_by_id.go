package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.GetByID")
	defer span.End()
	return q.signatures.GetByID(ctx, id)
}
