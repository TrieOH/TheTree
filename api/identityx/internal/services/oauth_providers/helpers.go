package oauth_providers

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

// requireProjectAdmin gates a write on the authenticated actor holding an
// admin or owner role in the provider's project. The caller is resolved
// by the Access-check module from the context identity.
func (o *Operations) requireProjectAdmin(ctx context.Context, projectID uuid.UUID) error {
	return o.authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
}
