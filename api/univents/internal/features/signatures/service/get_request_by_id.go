package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	return o.Requests.GetRequestByID(ctx, id)
}
