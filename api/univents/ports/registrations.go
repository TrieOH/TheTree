package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// RegistrationRepo is the registration read surface. The checkout feature owns
// the write side; readers (badges, certifications) consume confirmed rows.
type RegistrationRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Registration, error)
}
