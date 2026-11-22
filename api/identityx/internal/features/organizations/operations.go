// Package organizations is the organizations feature: its repo, handlers and
// operations. Build wires it.
package organizations

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"encoding/json"
	"errors"
	"lib/database"
	"lib/validator"
	"time"
	"uuid"
)

type Deps struct {
	Projects ports.ProjectRepo
	Actors   ports.ActorRepo
	Orgs     ports.OrganizationRepo
	// Keys owns the Key-lifecycle: org-created projects cross its Ensure
	// seam so they ship with keys instead of being token-broken until the
	// next boot.
	Keys  *keys.Manager
	Authz *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// projectInOrg returns the project when it belongs to the organization. An
// org route that names a project of another org (or a standalone one) gets
// not found: an org role reaches only the org's own projects.
func (o *Operations) projectInOrg(ctx context.Context, orgID, projectID uuid.UUID) (*models.Project, error) {
	project, err := o.Projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.OrganizationID == nil || *project.OrganizationID != orgID {
		return nil, errs.ErrProjectNotFound
	}
	return project, nil
}

func (o *Operations) AddMember(ctx context.Context, payload models.AddOrganizationMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return errs.ErrAddSelfToOrg
	}

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return errs.ErrAddOwnerToOrg
	}

	// Granting owner takes an owner; admins grant member and admin.
	minRole := models.OrganizationRoleAdmin
	if payload.Role == models.OrganizationRoleOwner {
		minRole = models.OrganizationRoleOwner
	}
	err = o.Authz.CheckOrg(ctx, payload.OrganizationID, minRole)
	if err != nil {
		return err
	}

	_, err = o.Actors.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	_, err = o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyOrgMember
	}

	newMember := models.OrganizationMember{
		ActorID:        actor.ID,
		OrganizationID: payload.OrganizationID,
		Role:           payload.Role,
	}

	return o.Orgs.AddMember(ctx, newMember)
}

func (o *Operations) AddProjectMember(ctx context.Context, payload models.AddOrgProjectMemberInput) error {
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

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return errs.ErrOrgOwnerAlreadyInProject
	}

	// Granting owner takes an org owner; org admins grant member and admin.
	minRole := models.OrganizationRoleAdmin
	if payload.Role == models.ProjectRoleOwner {
		minRole = models.OrganizationRoleOwner
	}
	err = o.Authz.CheckOrg(ctx, org.ID, minRole)
	if err != nil {
		return err
	}

	project, err := o.projectInOrg(ctx, org.ID, payload.ProjectID)
	if err != nil {
		return err
	}

	_, err = o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrOrgMemberAlreadyInProject
	}

	_, err = o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyProjectMember
	}

	return o.Projects.AddMember(ctx, models.ProjectMember{
		ProjectID: payload.ProjectID,
		ActorID:   actor.ID,
		Role:      payload.Role,
		Metadata:  json.RawMessage("{}"),
	})
}

func (o *Operations) Create(ctx context.Context, in models.CreateOrganizationInput) (*models.Organization, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org := &models.Organization{OwnerID: ident.Sub.ID, Name: in.Name, Slug: in.Slug}
	err = validator.Validate(org)
	if err != nil {
		return nil, err
	}

	var created *models.Organization
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Orgs.Create(ctx, *org)
		if err != nil {
			return err
		}

		owner := models.OrganizationMember{
			ActorID:        ident.Sub.ID,
			OrganizationID: created.ID,
			Role:           models.OrganizationRoleOwner,
		}

		return o.Orgs.AddMember(ctx, owner)
	})
	if err != nil {
		return nil, err
	}
	return created, err
}

func (o *Operations) CreateProject(ctx context.Context, in models.CreateOrgProjectInput) (*models.Project, error) {
	org, err := o.Orgs.GetByID(ctx, in.OrganizationID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckOrg(ctx, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		OrganizationID: &in.OrganizationID,
		OwnerID:        org.OwnerID,
		Name:           in.Name,
		BrandSlug:      in.BrandSlug,
		Domain:         in.Domain,
		Metadata:       json.RawMessage("{}"),
	}
	err = validator.Validate(project)
	if err != nil {
		return nil, err
	}

	// Project creation and key provisioning are one transaction: the
	// Key-lifecycle module provisions signing and encryption keys for the
	// new project here, so an org-created project is never token-broken
	// until the next boot.
	var created *models.Project
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Projects.Create(ctx, *project)
		if err != nil {
			return err
		}
		return o.Keys.Ensure(ctx, &created.ID)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (o *Operations) CreateProjectActor(ctx context.Context, payload models.CreateActorInput) (*models.Actor, error) {
	err := o.Authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleMember)
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

func (o *Operations) GetActorByID(ctx context.Context, id, projectID uuid.UUID) (*models.Actor, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
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

func (o *Operations) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	err = o.Authz.CheckOrg(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (o *Operations) ListOrgProjects(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	err := o.Authz.CheckOrg(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Projects.ListByOrganization(ctx, orgID)
}

func (o *Operations) ListOrgs(ctx context.Context) ([]models.Organization, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownOrgs, err := o.Orgs.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedOrgs, err := o.Orgs.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownOrgs, joinedOrgs...), nil
}

func (o *Operations) ListProjectActors(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Actors.List(ctx, projectID)
}

func (o *Operations) ListOrgProjectMembers(ctx context.Context, orgID, projectID uuid.UUID) ([]models.ProjectMember, error) {
	_, err := o.projectInOrg(ctx, orgID, projectID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	members, err := o.Projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := o.Orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	for _, m := range orgMembers {
		members = append(members, models.ProjectMember{
			ActorID:   m.ActorID,
			ProjectID: projectID,
			Role:      models.ProjectRole(m.Role),
		})
	}

	return members, nil
}

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveOrganizationMemberInput) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.Actors.GetByEmail(ctx, payload.ActorEmail, nil)
	if err != nil {
		return err
	}

	if ident.Sub.ID == actor.ID {
		return errs.ErrRemoveSelfFromOrg
	}

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}
	if actor.ID == org.OwnerID {
		return errs.ErrRemoveOrgOwner
	}

	err = o.Authz.CheckOrg(ctx, payload.OrganizationID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	member, err := o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotOrgMember
	}

	// Removing an admin or an owner takes an owner.
	if member.Role != models.OrganizationRoleMember {
		err = o.Authz.CheckOrg(ctx, org.ID, models.OrganizationRoleOwner)
		if err != nil {
			return err
		}
	}

	return o.Orgs.RemoveMember(ctx, actor.ID, payload.OrganizationID)
}

func (o *Operations) RemoveProjectMember(ctx context.Context, payload models.RemoveOrgProjectMemberInput) error {
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

	org, err := o.Orgs.GetByID(ctx, payload.OrganizationID)
	if err != nil {
		return err
	}

	if actor.ID == org.OwnerID {
		return errs.ErrRemoveOrgOwnerFromProject
	}

	err = o.Authz.CheckOrg(ctx, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	project, err := o.projectInOrg(ctx, org.ID, payload.ProjectID)
	if err != nil {
		return err
	}
	if actor.ID == project.OwnerID {
		return errs.ErrRemoveOrgOwnerFromProject
	}

	_, err = o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrRemoveOrgMemberFromProject
	}

	member, err := o.Projects.GetMember(ctx, actor.ID, project.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotProjectMember
	}

	// Removing a project admin or owner takes an org owner.
	if member.Role != models.ProjectRoleMember {
		err = o.Authz.CheckOrg(ctx, org.ID, models.OrganizationRoleOwner)
		if err != nil {
			return err
		}
	}

	return o.Projects.RemoveMember(ctx, actor.ID, payload.ProjectID)
}
