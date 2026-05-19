package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *QueryService) ListMembers(ctx context.Context, formID uuid.UUID) (members []models.FormMember, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.ListMembers")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var form *models.Form
	form, err = s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.NamespaceID != nil {
		var namespace *models.Namespace
		namespace, err = s.namespaces.GetByID(ctx, *form.NamespaceID)
		if err != nil {
			return nil, err
		}
		if sub.ID != namespace.OwnerID {
			_, err = s.namespaces.GetMember(ctx, sub.ID, namespace.ID)
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	} else {
		if sub.ID != form.OwnerID {
			_, err = s.forms.GetMember(ctx, sub.ID, formID)
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	members, err = s.forms.ListDirectMembers(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.NamespaceID != nil {
		var namespaceMembers []models.NamespaceMember
		namespaceMembers, err = s.namespaces.ListMembers(ctx, *form.NamespaceID)
		if err != nil {
			return nil, err
		}
		for _, m := range namespaceMembers {
			members = append(members, models.FormMember{
				UserID:  m.UserID,
				FormID:  form.ID,
				Role:    models.FormMemberRole(m.Role),
				AddedAt: m.AddedAt,
				AddedBy: m.AddedBy,
			})
		}
	}

	ownerMember := models.FormMember{
		UserID:  form.OwnerID,
		FormID:  form.ID,
		Role:    models.FormMemberRoleOwner,
		AddedAt: form.CreatedAt,
		AddedBy: form.OwnerID,
	}

	members = append(members, ownerMember)

	return members, nil
}
