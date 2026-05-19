package commands

import (
	"Informd/models"
	"context"
	"lib/authz"
	"time"

	"github.com/MintzyG/fun"
)

func (s *CommandService) AddMember(ctx context.Context, payload models.AddNamespaceMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.AddMember")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to namespaces")
	}

	var namespace *models.Namespace
	namespace, err = s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("owners can't be added to namespaces they own")
	}

	if sub.ID != namespace.OwnerID {
		var member *models.NamespaceMember
		member, err = s.namespaces.GetMember(ctx, sub.ID, payload.NamespaceID)
		if err != nil {
			return err
		}
		if member.Role != models.NamespaceMemberRoleAdmin {
			return fun.ErrForbidden("insufficient permissions")
		}
	}

	newMember := models.NamespaceMember{
		UserID:      payload.UserID,
		NamespaceID: payload.NamespaceID,
		Role:        payload.Role,
		AddedAt:     time.Now(),
		AddedBy:     sub.ID,
	}

	if err = s.namespaces.AddMember(ctx, newMember); err != nil {
		return err
	}
	return nil
}
