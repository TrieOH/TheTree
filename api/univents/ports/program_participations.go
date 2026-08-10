package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// ProgramParticipationRepo is the program_participations write side: created
// registered at checkout (split 7) attached to the required ticket item's
// registration, flipped by the webhook receiver (split 4) and the expiry
// worker (split 7).
type ProgramParticipationRepo interface {
	Create(ctx context.Context, toCreate *models.ProgramParticipation) (*models.ProgramParticipation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error)
}
