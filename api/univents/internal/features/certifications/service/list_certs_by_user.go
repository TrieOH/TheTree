package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListCertsByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	return o.Certs.ListByUser(ctx, userID)
}
