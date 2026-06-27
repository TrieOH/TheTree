package commands

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (s *Commands) AddMember(ctx context.Context, payload models.AddNamespaceMemberInput) error {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.AddMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to namespaces")
	}

	namespace, err := s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("owners can't be added to namespaces they own")
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
	if err == nil {
		return fun.ErrBadRequest("user is already a member of the namespace")
	}

	newMember := models.NamespaceMember{
		UserID:      payload.UserID,
		NamespaceID: payload.NamespaceID,
		Role:        payload.Role,
		AddedAt:     time.Now(),
		AddedBy:     ident.Sub.ID,
	}

	if err = s.namespaces.AddMember(ctx, newMember); err != nil {
		return err
	}
	return nil
}
