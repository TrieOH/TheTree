package authz

import (
	"Informd/ports"
	"context"

	libauthz "lib/authz"

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

func (s *Service) CheckForm(ctx context.Context, actorID, formID uuid.UUID, minRole libauthz.Role) error {
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

func (s *Service) CheckNamespace(ctx context.Context, actorID, namespaceID uuid.UUID, minRole libauthz.Role) error {
	role, err := s.namespaces.GetRole(ctx, actorID, namespaceID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
