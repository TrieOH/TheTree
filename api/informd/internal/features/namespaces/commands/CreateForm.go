package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *CommandService) CreateForm(ctx context.Context, title string, namespaceID uuid.UUID) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.CreateForm")
	defer span.End()

	var sub *authz.UserSubject
	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	var member *models.NamespaceMember
	if sub.ID != namespace.OwnerID {
		member, err = s.namespaces.GetMember(ctx, sub.ID, namespace.ID)
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
		if member.Role == models.NamespaceMemberRoleViewer {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	form, err := models.NewForm(&namespaceID, sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err := s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
