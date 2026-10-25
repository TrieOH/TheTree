package authz

import (
	"Informd/ports"
	"context"

	libauthz "lib/authz"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Service struct {
	forms      ports.FormsRepo
	namespaces ports.NamespaceRepo
}

func New(forms ports.FormsRepo, namespaces ports.NamespaceRepo) *Service {
	return &Service{forms: forms, namespaces: namespaces}
}

// Form resolves the caller from the identity the auth middleware wrote and
// checks they hold at least minRole on the form, directly or through its
// namespace. It returns the caller, so operations don't look it up again.
func (s *Service) Form(ctx context.Context, formID uuid.UUID, minRole libauthz.Role) (*idx.Identity, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return ident, s.checkForm(ctx, ident.Sub.ID, formID, minRole)
}

// Namespace resolves the caller like Form and checks they hold at least
// minRole on the namespace.
func (s *Service) Namespace(ctx context.Context, namespaceID uuid.UUID, minRole libauthz.Role) (*idx.Identity, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return ident, s.checkNamespace(ctx, ident.Sub.ID, namespaceID, minRole)
}

func (s *Service) checkForm(ctx context.Context, actorID, formID uuid.UUID, minRole libauthz.Role) error {
	form, err := s.forms.GetByID(ctx, formID)
	if err != nil {
		return err
	}

	formRole, err := s.forms.GetRole(ctx, actorID, formID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil && libauthz.Min(formRole, minRole) == nil {
		return nil
	}

	if form.NamespaceID == nil {
		if err != nil {
			return libauthz.ForbiddenIfNotFound(err)
		}
		return libauthz.Min(formRole, minRole)
	}

	namespaceRole, err := s.namespaces.GetRole(ctx, actorID, *form.NamespaceID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(namespaceRole, minRole)
}

func (s *Service) checkNamespace(ctx context.Context, actorID, namespaceID uuid.UUID, minRole libauthz.Role) error {
	role, err := s.namespaces.GetRole(ctx, actorID, namespaceID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
