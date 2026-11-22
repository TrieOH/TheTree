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

	libauthz "lib/authz"
)

type Deps struct {
	Actors   ports.ActorRepo
	Projects ports.ProjectRepo
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// requireActorReadAccess gates reads of the project's actors to the
// project's own backends and to its members. A caller scoped to the project
// passes when it is a backend: an API key, or a service or machine actor.
// A human signed in to the project is one of its end users and never sees
// the others. A platform-level client (nil project_id) passes only as a
// member, never by virtue of being platform-level. Never public.
func (o *Operations) requireActorReadAccess(ctx context.Context, ident *models.Identity, projectID uuid.UUID) error {
	if ident.Sub.ProjectID != nil && *ident.Sub.ProjectID == projectID {
		if ident.Cred.Type == models.APIKeyCredentialType || ident.Sub.Type != models.HumanActorType {
			return nil
		}
		return libauthz.ErrInsufficientRole
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

	// An end user may read their own actor; everything else needs read access.
	if ident.Sub.ID != id {
		err = o.requireActorReadAccess(ctx, ident, projectID)
		if err != nil {
			return nil, err
		}
	}

	actor, err := o.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actor.ProjectID == nil || *actor.ProjectID != projectID {
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
