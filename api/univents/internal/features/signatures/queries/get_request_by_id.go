package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.GetRequestByID")
	defer span.End()
	return q.requests.GetRequestByID(ctx, id)
}
