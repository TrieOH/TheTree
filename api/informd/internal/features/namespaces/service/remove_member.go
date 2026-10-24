package service

import (
	"context"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.RemoveMember")
	defer span.End()

	namespace, err := o.Namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Namespace(ctx, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the namespace")
	}

	return o.Namespaces.RemoveMember(ctx, payload.UserID, payload.NamespaceID)
}
