package namespaces

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (o *Operations) RemoveFormMember(ctx context.Context, payload models.RemoveNamespaceFormMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.RemoveFormMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't remove themselves from forms")
	}

	namespace, err := o.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("cannot remove owner of the namespace from form")
	}

	err = o.authz.CheckNamespace(ctx, ident.Sub.ID, namespace.ID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("cannot remove namespace member from form")
	}

	form, err := o.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	_, err = o.forms.GetMember(ctx, payload.UserID, form.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the form")
	}

	return o.forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
