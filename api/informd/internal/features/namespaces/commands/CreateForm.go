package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Commands) CreateForm(ctx context.Context, title string, namespaceID uuid.UUID) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.CreateForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	var member *models.NamespaceMember
	if ident.Sub.ID != namespace.OwnerID {
		member, err = s.namespaces.GetMember(ctx, ident.Sub.ID, namespace.ID)
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
		if member.Role == models.NamespaceMemberRoleViewer {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	form, err := models.NewForm(&namespaceID, namespace.OwnerID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err := s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
