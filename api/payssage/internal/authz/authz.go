package authz

import (
	"context"

	libauthz "lib/authz"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type service struct {
	orgs    ports.OrganizationRepo
	wallets ports.WalletRepo
}

func New(orgs ports.OrganizationRepo, wallets ports.WalletRepo) *service {
	return &service{orgs: orgs, wallets: wallets}
}

var Service *service

func (s *service) CheckOrg(ctx context.Context, actorID, orgID uuid.UUID, minRole libauthz.Role) error {
	role, err := s.orgs.GetRole(ctx, actorID, orgID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}

func (s *service) CheckWalletAccess(ctx context.Context, actorID, walletID uuid.UUID, minRole libauthz.Role) error {
	wallet, err := s.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	role, err := s.wallets.GetRole(ctx, actorID, walletID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil && libauthz.Min(role, minRole) == nil {
		return nil
	}

	if wallet.OrganizationID == nil {
		if err != nil {
			return libauthz.ForbiddenIfNotFound(err)
		}
		return libauthz.Min(role, minRole)
	}

	orgRole, err := s.orgs.GetRole(ctx, actorID, *wallet.OrganizationID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(orgRole, minRole)
}
