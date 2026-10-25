package authz

import (
	"context"

	libauthz "lib/authz"
	"payssage/ports"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Service struct {
	orgs    ports.OrganizationRepo
	wallets ports.WalletRepo
}

func New(orgs ports.OrganizationRepo, wallets ports.WalletRepo) *Service {
	return &Service{orgs: orgs, wallets: wallets}
}

// Org resolves the caller from the identity the auth middleware wrote and
// checks they hold at least minRole on the organization. It returns the
// caller, so operations don't look it up again.
func (s *Service) Org(ctx context.Context, orgID uuid.UUID, minRole libauthz.Role) (*idx.Identity, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return ident, s.checkOrg(ctx, ident.Sub.ID, orgID, minRole)
}

// Wallet resolves the caller like Org and checks they hold at least
// minRole on the wallet, directly or through its organization.
func (s *Service) Wallet(ctx context.Context, walletID uuid.UUID, minRole libauthz.Role) (*idx.Identity, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return ident, s.checkWallet(ctx, ident.Sub.ID, walletID, minRole)
}

func (s *Service) checkOrg(ctx context.Context, actorID, orgID uuid.UUID, minRole libauthz.Role) error {
	role, err := s.orgs.GetRole(ctx, actorID, orgID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}

func (s *Service) checkWallet(ctx context.Context, actorID, walletID uuid.UUID, minRole libauthz.Role) error {
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
