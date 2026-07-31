package queries

import (
	"context"
	"payssage/models"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Queries struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewQueries(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		wallets: wallets,
		orgs:    orgs,
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
