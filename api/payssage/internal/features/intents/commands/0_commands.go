package commands

import (
	"context"
	"fmt"
	"lib/database"
	"payssage/models"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	intents    ports.IntentRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewCommands(
	intents ports.IntentRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		intents:    intents,
		wallets:    wallets,
		orgs:       orgs,
		collectors: collectors,
		sellers:    sellers,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}

func (c *Commands) checkRole(ctx context.Context, org *models.Organization, subID uuid.UUID, minRole models.OrganizationRole) error {
	if org == nil {
		return fun.ErrForbidden("insufficient permissions")
	}

	if org.OwnerID == subID {
		return nil
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

func (c *Commands) checkAdminAccess(ctx context.Context, walletID, subID uuid.UUID) error {
	wallet, err := c.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	// personal wallet: owner is admin
	if wallet.OrganizationID == nil {
		if wallet.OwnerID != subID {
			return fun.ErrForbidden("insufficient permissions")
		}
		return nil
	}

	// org wallet: check org admin role
	org, err := c.orgs.GetByID(ctx, *wallet.OrganizationID)
	if err != nil {
		return err
	}

	if err := c.checkRole(ctx, org, subID, models.OrganizationRoleAdmin); err != nil {
		return err
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
