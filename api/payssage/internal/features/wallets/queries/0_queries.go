package queries

import (
	"context"
	"lib/database"
	"payssage/models"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Queries struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	tracer  trace.Tracer
	tx      database.TxRunner
}

func NewQueries(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		wallets: wallets,
		orgs:    orgs,
		tracer:  tracer,
		tx:      tx,
	}
}

func (q *Queries) checkRole(ctx context.Context, org *models.Organization, subID uuid.UUID, minRole models.OrganizationRole) error {
	if org == nil {
		return fun.ErrForbidden("insufficient permissions")
	}

	if org.OwnerID == subID {
		return nil // owner always passes
	}

	member, err := q.orgs.GetMember(ctx, subID, org.ID)
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
