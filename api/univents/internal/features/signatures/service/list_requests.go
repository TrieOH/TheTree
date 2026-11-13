package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error) {
	return o.Requests.ListRequestsByEdition(ctx, editionID)
}
