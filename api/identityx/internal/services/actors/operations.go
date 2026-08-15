package actors

import (
	"context"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/errx"

	"github.com/google/uuid"
)

type Operations struct {
	actors   ports.ActorRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func NewOperations(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		actors:   actors,
		projects: projects,
		authz:    authz,
	})
}

// requireActorReadAccess gates actor reads to callers with a relationship
// to the project: any user of the project (project_id matches — svc accounts
// included, no membership row needed) or any member of the project. A
// platform-level client (nil project_id) passes only as a member, never by
// virtue of being platform-level. Never public.
func (o *Operations) requireActorReadAccess(ctx context.Context, ident *models.Identity, projectID uuid.UUID) error {
	if ident.Sub.ProjectID != nil && *ident.Sub.ProjectID == projectID {
		return nil
	}
	return o.authz.CheckProject(ctx, ident.Sub.ID, projectID, models.ProjectRoleMember)
}
