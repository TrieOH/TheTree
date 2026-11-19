// Package actors is the actors feature: its repo, handlers and operations.
// Build wires it.
package actors

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"time"
	"uuid"
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

func (o *Operations) Create(ctx context.Context, payload models.CreateActorInput) (*models.Actor, error) {
	err := o.Authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	if payload.Type != models.HumanActorType && payload.AuthMethod == models.PasswordAuthMethod {
		return nil, errs.ErrPasswordOnlyForHumans
	}

	actor := models.Actor{
		ProjectID:  payload.ProjectID,
		AuthMethod: payload.AuthMethod,
		VerifiedAt: new(time.Now()),
		Email:      payload.Email,
		Type:       payload.Type,
	}

	return o.Actors.Register(ctx, actor)
}

func (o *Operations) GetByEmail(ctx context.Context, email string, projectID uuid.UUID) (*models.Actor, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.requireActorReadAccess(ctx, ident, projectID)
	if err != nil {
		return nil, err
	}

	actor, err := o.Actors.GetByEmail(ctx, email, &projectID)
	if err != nil {
		return nil, err
	}

	return actor, nil
}

func (o *Operations) GetByID(ctx context.Context, id, projectID uuid.UUID) (*models.Actor, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.requireActorReadAccess(ctx, ident, projectID)
	if err != nil {
		return nil, err
	}

	actor, err := o.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actor.ProjectID != nil && *actor.ProjectID != projectID {
		return nil, errs.ErrActorOutsideProject
	}

	return actor, nil
}

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.requireActorReadAccess(ctx, ident, projectID)
	if err != nil {
		return nil, err
	}

	return o.Actors.List(ctx, projectID)
}
