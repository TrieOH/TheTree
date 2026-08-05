package namespaces

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (o *Operations) AddFormMember(ctx context.Context, payload models.AddNamespaceFormMemberInput) error {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.AddFormMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to forms")
	}

	namespace, err := o.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("owner of the namespace is already a member of the form")
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
		return fun.ErrBadRequest("namespace member is already a member of the form")
	}

	form, err := o.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	_, err = o.forms.GetMember(ctx, payload.UserID, form.ID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return fun.ErrBadRequest("user is already a member of the form")
	}

	newMember := models.FormMember{
		UserID:  payload.UserID,
		FormID:  form.ID,
		Role:    payload.Role,
		AddedAt: time.Now(),
		AddedBy: ident.Sub.ID,
	}

	return o.forms.AddMember(ctx, newMember)
}
