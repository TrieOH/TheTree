package ports

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

type ResponderRepo interface {
	Create(ctx context.Context, toCreate models.Responder) (*models.Responder, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Responder, error)
	GetByEmail(ctx context.Context, email string) (*models.Responder, error)
	GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Responder, error)
}
