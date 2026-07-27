package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type ProgramRepo interface {
	Create(ctx context.Context, toCreate *models.Program) (*models.Program, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Program, error)
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error)
	Patch(ctx context.Context, id uuid.UUID, program *models.Program) (*models.Program, error)
	Delete(ctx context.Context, id uuid.UUID) (*models.Program, error)
}

type ProgramOccurrenceRepo interface {
	CreateOccurrence(ctx context.Context, toCreate *models.ProgramOccurrence) (*models.ProgramOccurrence, error)
	GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error)
	ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error)
	ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error)
	PatchOccurrence(ctx context.Context, id uuid.UUID, occurrence *models.ProgramOccurrence) (*models.ProgramOccurrence, error)
	DeleteOccurrence(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error)
}
