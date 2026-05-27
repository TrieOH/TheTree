package commands

import (
	models2 "IdentityX/models"
	"Informd/models"
	"context"

	"github.com/MintzyG/fun"
)

func (s *CommandService) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) error {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.RemoveMember")
	defer span.End()

	sub, err := models2.RequireSubject(ctx)
	if err != nil {
		return err
	}

	namespace, err := s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if sub.ID != namespace.OwnerID {
		member, err := s.namespaces.GetMember(ctx, sub.ID, payload.NamespaceID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err != nil {
			return fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.NamespaceMemberRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	_, err = s.namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the namespace")
	}

	return s.namespaces.RemoveMember(ctx, payload.UserID, payload.NamespaceID)
}
