// Package wallets is the wallets feature: its repo, handlers and operations.
// Build wires it.
package wallets

import (
	"context"
	"payssage/errs"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	idx "sdk/identityx"
	"uuid"
)

type Deps struct {
	Wallets ports.WalletRepo
	Orgs    ports.OrganizationRepo
	Authz   *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet.OrganizationID != nil {
		org, err := o.Orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return err
		}
		_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return err
		}
	} else if wallet.OwnerID != ident.Sub.ID {
		return errs.ErrNotResourceOwner
	}

	return o.Wallets.BindCollector(ctx, walletID, collectorID)
}

func (o *Operations) Create(ctx context.Context, payload models.CreateWalletInput) (*models.Wallet, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var org *models.Organization
	if payload.OrganizationID != nil {
		org, err = o.Orgs.GetByID(ctx, *payload.OrganizationID)
		if err != nil {
			return nil, err
		}
		_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return nil, err
		}
	}

	ownerID := ident.Sub.ID
	var orgID *uuid.UUID
	if org != nil {
		ownerID = org.OwnerID
		orgID = &org.ID
	}

	wallet := models.Wallet{
		OwnerID:        ownerID,
		OrganizationID: orgID,
		Name:           payload.Name,
		Sandbox:        false,
		FeeBps:         0,
	}

	return o.Wallets.Create(ctx, wallet)
}

func (o *Operations) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (o *Operations) ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error) {
	org, err := o.Orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Wallets.ListFromOrg(ctx, org.ID)
}

func (o *Operations) List(ctx context.Context) ([]models.Wallet, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.Wallets.List(ctx, ident.Sub.ID)
}

func (o *Operations) SetFeeBPS(ctx context.Context, payload models.SetFeeBPSInput) error {
	wallet, err := o.Wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return o.Wallets.SetFeeBPS(ctx, wallet.ID, payload.FeeBps)
}

func (o *Operations) SetSandbox(ctx context.Context, payload models.SetSandboxInput) error {
	wallet, err := o.Wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return o.Wallets.SetSandboxState(ctx, wallet.ID, payload.Sandbox)
}

func (o *Operations) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet.OrganizationID != nil {
		org, err := o.Orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return err
		}
		_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return err
		}
	} else if wallet.OwnerID != ident.Sub.ID {
		return errs.ErrNotResourceOwner
	}

	return o.Wallets.UnbindCollector(ctx, walletID)
}
