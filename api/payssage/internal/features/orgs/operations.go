// Package orgs is the orgs feature: its repo, handlers and operations. Build
// wires it.
package orgs

import (
	"context"
	"errors"
	"lib/database"
	"payssage/errs"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	idx "sdk/identityx"
	"uuid"
)

type Deps struct {
	Orgs  ports.OrganizationRepo
	IDX   *idx.Client
	Authz *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) AddMember(ctx context.Context, payload models.AddOrganizationMemberInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, payload.ActorEmail)
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

	_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
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
		MemberID:       actor.ID,
		OrganizationID: payload.OrganizationID,
		Role:           payload.Role,
	}

	return o.Orgs.AddMember(ctx, newMember)
}

func (o *Operations) Create(ctx context.Context, in models.CreateOrganizationInput) (*models.Organization, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org := models.Organization{
		OwnerID: ident.Sub.ID,
		Name:    in.Name,
		Slug:    in.Slug,
	}

	var created *models.Organization
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Orgs.Create(ctx, org)
		if err != nil {
			return err
		}

		owner := models.OrganizationMember{
			MemberID:       ident.Sub.ID,
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

func (o *Operations) GetMemberByEmail(ctx context.Context, email string, orgID uuid.UUID) (*models.OrganizationMember, error) {
	_, err := o.Authz.Org(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	member, err := o.Orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (o *Operations) GetMemberByID(ctx context.Context, id, orgID uuid.UUID) (*models.OrganizationMember, error) {
	_, err := o.Authz.Org(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := o.IDX.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	member, err := o.Orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (o *Operations) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	_, err = o.Authz.Org(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (o *Operations) ListOrgs(ctx context.Context) ([]models.Organization, error) {
	ident, err := idx.RequireIdentity(ctx)
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

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveOrganizationMemberInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, payload.ActorEmail)
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

	_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Orgs.GetMember(ctx, actor.ID, org.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotOrgMember
	}

	return o.Orgs.RemoveMember(ctx, actor.ID, payload.OrganizationID)
}
