package commands

import (
	"context"

	"Informd/models"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) EditSelectConfigNamespaced(ctx context.Context, formID, namespaceID uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := s.tracer.Start(ctx, "FieldService.EditSelectConfigNamespaced")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	namespaceMember, err := s.namespaces.GetMember(ctx, sub.ID, namespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}
	if fun.Is(err, fun.CodeNotFound) {
		if namespaceMember.Role == models.NamespaceMemberRoleViewer {
			member, err := s.forms.GetMember(ctx, sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
			if member.Role == models.FormMemberRoleViewer {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	return s.fields.UpdateSelectConfig(ctx, payload)
}
