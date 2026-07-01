package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type CapabilityRepo interface {
	Create(ctx context.Context, toCreate models.Capability) (*models.Capability, error)
	List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error)
	ValidateCapabilities(ctx context.Context, projectID *uuid.UUID, capabilities []uuid.UUID) (bool, error)
	AssignToApiKey(ctx context.Context, apiKeyID uuid.UUID, capabilityIDs []uuid.UUID, assignedBy uuid.UUID) error
	ListByApiKeyPrefix(ctx context.Context, prefix string) ([]models.Capability, error)
}
