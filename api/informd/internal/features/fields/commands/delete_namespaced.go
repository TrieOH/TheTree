package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) DeleteNamespaced(ctx context.Context, namespaceID, formID, fieldID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "FieldService.DeleteNamespaced")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	namespaceMember, err := s.namespaces.GetMember(ctx, sub.ID, namespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if fun.Is(err, fun.CodeNotFound) {
		if namespaceMember.Role == models.NamespaceMemberRoleViewer {
			member, err := s.forms.GetMember(ctx, sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return err
			}
			if err != nil {
				return fun.ErrForbidden("insufficient permissions")
			}
			if member.Role == models.FormMemberRoleViewer {
				return fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	return s.fields.Delete(ctx, fieldID)
}
