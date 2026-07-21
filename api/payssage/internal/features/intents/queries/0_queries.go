package queries

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

type Queries struct {
	intents ports.IntentRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	logger  *zap.Logger
	tracer  trace.Tracer
	tx      database.TxRunner
}

func NewQueries(
	intents ports.IntentRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		intents: intents,
		wallets: wallets,
		orgs:    orgs,
		logger:  logger,
		tracer:  tracer,
		tx:      tx,
	}
}

func (q *Queries) checkWalletAccess(ctx context.Context, walletID, subID uuid.UUID) error {
	wallet, err := q.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet.OwnerID == subID {
		return nil
	}

	if wallet.OrganizationID != nil {
		org, err := q.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return err
		}
		if org.OwnerID == subID {
			return nil
		}
		_, err = q.orgs.GetMember(ctx, subID, org.ID)
		if err != nil {
			return fmt.Errorf("insufficient permissions: %w", err)
		}
		return nil
	}

	return fun.ErrForbidden("insufficient permissions")
}
