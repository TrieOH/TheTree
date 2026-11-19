// Package projects is the projects feature: its repo, handlers and
// operations. Build wires it.
package projects

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"errors"
	"lib/database"
	"uuid"
)

type Deps struct {
	// Keys owns the Key-lifecycle: project creation crosses its Ensure
	// seam so every project ships with signing and encryption keys.
	Keys     *keys.Manager
	Projects ports.ProjectRepo
	Actors   ports.ActorRepo
	Authz    *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) AddMember(ctx context.Context, payload models.AddProjectMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return errs.ErrAddSelfToProject
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	if actor.ID == project.OwnerID {
		return errs.ErrAddOwnerToProject
	}

	err = o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Actors.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	_, err = o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyProjectMember
	}

	newMember, err := models.NewProjectMember(payload.ProjectID, actor.ID, payload.Role)
	if err != nil {
		return err
	}

	return o.Projects.AddMember(ctx, *newMember)
}

func (o *Operations) Create(ctx context.Context, in models.CreateProjectInput) (*models.Project, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := models.NewProject(ident.Sub.ID, in.BrandSlug, in.Name, in.Domain, nil)
	if err != nil {
		return nil, err
	}

	var created *models.Project
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Projects.Create(ctx, *project)
		if err != nil {
			return err
		}

		// The Key-lifecycle module provisions the project's signing and
		// encryption keys inside the same tx; an org-created project (the
		// organizations feature) crosses the same seam, so every project
		// ships with keys regardless of which path created it.
		err = o.Keys.Ensure(ctx, &created.ID)
		if err != nil {
			return err
		}

		member, err := models.NewProjectMember(created.ID, ident.Sub.ID, models.ProjectRoleOwner)
		if err != nil {
			return err
		}

		return o.Projects.AddMember(ctx, *member)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (o *Operations) ListProjects(ctx context.Context) ([]models.Project, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownProjects, err := o.Projects.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedProjects, err := o.Projects.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownProjects, joinedProjects...), nil
}

func (o *Operations) ListMembers(ctx context.Context, projectID uuid.UUID) (members []models.ProjectMember, err error) {
	err = o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveProjectMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return errs.ErrRemoveSelfFromProject
	}

	project, err := o.Projects.GetByID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}
	if actor.ID == project.OwnerID {
		return errs.ErrRemoveProjectOwner
	}

	err = o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotProjectMember
	}

	return o.Projects.RemoveMember(ctx, actor.ID, payload.ProjectID)
}
