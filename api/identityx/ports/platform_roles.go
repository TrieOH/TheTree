package ports

import (
	"IdentityX/models"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type PlatformRolesRepo interface {
	Give(ctx context.Context, actorID uuid.UUID, role models.PlatformRole, metadata *json.RawMessage) (*models.PlatformRoleRelation, error)
}
