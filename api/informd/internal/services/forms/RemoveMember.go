package forms

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveFormMemberInput) (err error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.RemoveMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't remove themselves from forms")
	}

	var form *models.Form
	form, err = o.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return fun.ErrBadRequest("cannot remove owner of the form")
	}
	err = o.authz.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.forms.GetMember(ctx, payload.UserID, form.ID)
	if !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrBadRequest("user is not a member of the form")
	}

	return o.forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
