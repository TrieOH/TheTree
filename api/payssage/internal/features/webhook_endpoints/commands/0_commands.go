package commands

import (
	"context"
	"fmt"
	"lib/database"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	endpoints ports.WebhookEndpointRepo
	wallets   ports.WalletRepo
	orgs      ports.OrganizationRepo
	logger    *zap.Logger
	tracer    trace.Tracer
	tx        database.TxRunner
}

func NewCommands(
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		endpoints: endpoints,
		wallets:   wallets,
		orgs:      orgs,
		logger:    logger,
		tracer:    tracer,
		tx:        tx,
	}
}

func (c *Commands) checkWalletAccess(ctx context.Context, walletID, subID uuid.UUID) error {
	wallet, err := c.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}
	if wallet.OwnerID == subID {
		return nil
	}
	if wallet.OrganizationID != nil {
		org, err := c.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return err
		}
		if org.OwnerID == subID {
			return nil
		}
		_, err = c.orgs.GetMember(ctx, subID, org.ID)
		if err != nil {
			return fmt.Errorf("insufficient permissions: %w", err)
		}
		return nil
	}
	return fun.ErrForbidden("insufficient permissions")
}
