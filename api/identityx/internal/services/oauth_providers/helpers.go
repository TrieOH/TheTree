package oauth_providers

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

// requireProjectAdmin gates a write on the authenticated actor holding an
// admin or owner role in the provider's project.
func (o *Operations) requireProjectAdmin(ctx context.Context, projectID uuid.UUID) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}
	return o.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleAdmin)
}
