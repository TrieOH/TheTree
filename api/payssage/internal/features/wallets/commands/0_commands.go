package commands

import (
	"context"
	"fmt"
	"payssage/models"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Commands struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewCommands(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Commands {
	return &Commands{
		wallets: wallets,
		orgs:    orgs,
	}
}

//nolint:unparam // TODO: extract into shared authzChecker service
func (c *Commands) checkRole(ctx context.Context, org *models.Organization, subID uuid.UUID, minRole models.OrganizationRole) error {
	if org == nil {
		return fun.ErrForbidden("insufficient permissions")
	}

	if org.OwnerID == subID {
		return nil // owner always passes
	}

	member, err := c.orgs.GetMember(ctx, subID, org.ID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return fun.ErrForbidden("insufficient permissions")
		}
		return err
	}

	if !member.Role.AtLeast(minRole) {
		return fun.ErrForbidden("insufficient permissions")
	}

	return nil
}

func (c *Commands) checkWalletAccess(wallet *models.Wallet, subID uuid.UUID, org ...*models.Organization) error {
	if wallet == nil {
		return fun.ErrForbidden("insufficient permissions")
	}

	if len(org) > 1 {
		return fmt.Errorf("checkWalletAccess: expected at most one org, got %d", len(org))
	}

	if len(org) == 1 && org[0] != nil {
		if wallet.OrganizationID == nil || *wallet.OrganizationID != org[0].ID {
			return fun.ErrForbidden("insufficient permissions")
		}
		return nil
	}

	if wallet.OwnerID != subID {
		return fun.ErrForbidden("insufficient permissions")
	}

	return nil
}
