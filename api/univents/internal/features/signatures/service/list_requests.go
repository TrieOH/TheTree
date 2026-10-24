package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.ListRequestsByEdition")
	defer span.End()
	return o.Requests.ListRequestsByEdition(ctx, editionID)
}
