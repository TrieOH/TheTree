package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
)

func (s *CommandService) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.RemoveMember")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var namespace *models.Namespace
	namespace, err = s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if sub.ID != namespace.OwnerID {
		var member *models.NamespaceMember
		member, err = s.namespaces.GetMember(ctx, sub.ID, payload.NamespaceID)
		if err != nil {
			return err
		}
		if member.Role != models.NamespaceMemberRoleAdmin {
			return fun.ErrForbidden("insufficient permission")
		}
	}

	if err = s.namespaces.RemoveMember(ctx, payload.UserID, payload.NamespaceID); err != nil {
		return err
	}
	return nil
}
