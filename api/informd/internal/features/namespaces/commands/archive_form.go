package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *CommandService) ArchiveForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.ArchiveForm")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	form, err := s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusClosed {
		return nil, fun.ErrBadRequest("cannot archive a form not on closed")
	}

	if sub.ID != namespace.OwnerID {
		member, err := s.namespaces.GetMember(ctx, sub.ID, namespaceID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.NamespaceMemberRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.forms.Close(ctx, form.ID)
}
