package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.GetRequestByID")
	defer span.End()
	return o.Requests.GetRequestByID(ctx, id)
}
