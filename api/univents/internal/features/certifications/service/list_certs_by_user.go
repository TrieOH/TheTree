package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListCertsByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	return o.Certs.ListByUser(ctx, userID)
}
