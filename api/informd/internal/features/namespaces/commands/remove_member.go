package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (s *Commands) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) error {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	namespace, err := s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if ident.Sub.ID != namespace.OwnerID {
		member, err := s.namespaces.GetMember(ctx, ident.Sub.ID, payload.NamespaceID)
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
