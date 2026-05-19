package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
)

func (s *CommandService) RemoveMember(ctx context.Context, payload models.RemoveFormMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.RemoveMember")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't remove themselves from forms")
	}

	var form *models.Form
	if form, err = s.forms.GetByID(ctx, payload.FormID); err != nil {
		return err
	}

	var namespace *models.Namespace
	if form.NamespaceID != nil {
		if namespace, err = s.namespaces.GetByID(ctx, *form.NamespaceID); err != nil {
			return err
		}
		if payload.UserID == namespace.OwnerID {
			return fun.ErrBadRequest("cannot remove owner of the namespace from form")
		}
		if sub.ID != namespace.OwnerID {
			if err = s.isFormAdmin(ctx, sub.ID, form.NamespaceID, payload.FormID); err != nil {
				return err
			}
		}
		_, err = s.namespaces.GetMember(ctx, payload.UserID, namespace.ID)
		if err == nil {
			return fun.ErrBadRequest("cannot remove namespace member from a namespace form")
		}
		if fun.Is(err, fun.CodeNotFound) {
			_, err = s.forms.GetMember(ctx, payload.UserID, payload.FormID)
			if fun.Is(err, fun.CodeNotFound) {
				return fun.ErrBadRequest("user already removed from form")
			}
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	} else {
		if payload.UserID == form.OwnerID {
			return fun.ErrBadRequest("cannot remove owner of the form")
		}
		if sub.ID != form.OwnerID {
			if err = s.isFormAdmin(ctx, sub.ID, form.NamespaceID, payload.FormID); err != nil {
				return err
			}
		}
		_, err = s.forms.GetMember(ctx, payload.UserID, payload.FormID)
		if fun.Is(err, fun.CodeNotFound) {
			return fun.ErrBadRequest("user already removed from form")
		}
		if err != nil {
			return err
		}
	}

	return s.forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
