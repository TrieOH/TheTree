package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// FIXME split this into two, namespaced and not this one should be only direct, namespaced should be handled by namespace feature

func (s *CommandService) Create(ctx context.Context, title string, namespaceID *uuid.UUID) (created *models.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	ownerID := sub.ID
	var namespace *models.Namespace
	if namespaceID != nil {
		namespace, err = s.namespaces.GetByID(ctx, *namespaceID)
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

		ownerID = namespace.OwnerID
	}

	var form *models.Form
	form, err = models.NewForm(namespaceID, ownerID, title)
	if err != nil {
		return nil, err
	}

	created, err = s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
