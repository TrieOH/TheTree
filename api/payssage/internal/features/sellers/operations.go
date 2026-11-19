// Package sellers is the sellers feature: its repo, handlers and operations.
// Build wires it.
package sellers

import (
	"context"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	"uuid"
)

type Deps struct {
	Sellers ports.SellerRepo
	Wallets ports.WalletRepo
	Orgs    ports.OrganizationRepo
	Authz   *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Seller, error) {
	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Wallet(ctx, wallet.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Sellers.ListByWallet(ctx, wallet.ID)
}
