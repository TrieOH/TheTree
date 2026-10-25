package service

import (
	"context"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/google/uuid"
)

type Deps struct {
	Actors   ports.ActorRepo
	Projects ports.ProjectRepo
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// requireActorReadAccess gates actor reads to callers with a relationship
// to the project: any user of the project (project_id matches — svc accounts
// included, no membership row needed) or any member of the project. A
// platform-level client (nil project_id) passes only as a member, never by
// virtue of being platform-level. Never public.
func (o *Operations) requireActorReadAccess(ctx context.Context, ident *models.Identity, projectID uuid.UUID) error {
	if ident.Sub.ProjectID != nil && *ident.Sub.ProjectID == projectID {
		return nil
	}
	return o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
}
