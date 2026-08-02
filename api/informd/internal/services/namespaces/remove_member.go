package namespaces

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	namespace, err := o.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	err = o.authz.CheckNamespace(ctx, ident.Sub.ID, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the namespace")
	}

	return o.namespaces.RemoveMember(ctx, payload.UserID, payload.NamespaceID)
}
