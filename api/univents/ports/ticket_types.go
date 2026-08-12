package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type TicketTypeRepo interface {
	Create(ctx context.Context, toCreate *models.TicketType) (*models.TicketType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.TicketType, error)
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error)
	Patch(ctx context.Context, id uuid.UUID, ticketType *models.TicketType) (*models.TicketType, error)
}
