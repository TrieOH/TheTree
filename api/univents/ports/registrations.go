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

	// Create + UpdateStatus are the checkout (split 7) / webhook receiver
	// (split 4) write side: pending rows are materialized at checkout and
	// flipped on approve/cancel/expire.
	Create(ctx context.Context, toCreate *models.Registration) (*models.Registration, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.RegistrationStatus, reason *string) (*models.Registration, error)
}
