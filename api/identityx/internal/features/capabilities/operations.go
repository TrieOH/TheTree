// Package capabilities is the capabilities feature: its repo, handlers and
// operations. Build wires it.
package capabilities

import (
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"uuid"
)

type Deps struct {
	Actors       ports.ActorRepo
	Capabilities ports.CapabilityRepo
	Projects     ports.ProjectRepo
	Authz        *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) Create(ctx context.Context, payload models.CreateCapabilityInput) (*models.Capability, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	capability := models.Capability{
		ProjectID: payload.ProjectID,
		Resource:  payload.Resource,
		Action:    payload.Action,
		CreatedBy: ident.Sub.ID,
	}

	return o.Capabilities.Create(ctx, capability)
}

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Capabilities.List(ctx, projectID)
}
