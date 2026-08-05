package namespaces

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddNamespaceMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.AddMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to namespaces")
	}

	namespace, err := o.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("owners can't be added to namespaces they own")
	}

	err = o.authz.CheckNamespace(ctx, ident.Sub.ID, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.namespaces.GetMember(ctx, payload.UserID, namespace.ID)
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

	err = o.namespaces.AddMember(ctx, newMember)
	if err != nil {
		return err
	}
	return nil
}
