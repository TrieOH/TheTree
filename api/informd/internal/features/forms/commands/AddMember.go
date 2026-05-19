package commands

import (
	"Informd/models"
	"context"
	"lib/authz"
	"time"

	"github.com/MintzyG/fun"
)

func (s *CommandService) AddMember(ctx context.Context, payload models.AddFormMemberInput) (err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.AddMember")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to forms")
	}

	var form *models.Form
	form, err = s.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if form.NamespaceID != nil {
		var namespace *models.Namespace
		if namespace, err = s.namespaces.GetByID(ctx, *form.NamespaceID); err != nil {
			return err
		}
		if payload.UserID == namespace.OwnerID {
			return fun.ErrBadRequest("owner of the namespace is already a member of the form")
		}
		if sub.ID != namespace.OwnerID {
			if err = s.isFormAdmin(ctx, sub.ID, form.NamespaceID, payload.FormID); err != nil {
				return err
			}
		}
		_, err = s.namespaces.GetMember(ctx, payload.UserID, namespace.ID)
		if err == nil {
			return fun.ErrBadRequest("namespace member is already a member of the form")
		}
		if fun.Is(err, fun.CodeNotFound) {
			_, err = s.forms.GetMember(ctx, payload.UserID, payload.FormID)
			if err == nil {
				return fun.ErrBadRequest("member is already a member of the form")
			}
			if !fun.Is(err, fun.CodeNotFound) {
				return err
			}
		} else if err != nil {
			return err
		}
	} else {
		if payload.UserID == form.OwnerID {
			return fun.ErrBadRequest("owner of the form is already a member of the form")
		}
		if sub.ID != form.OwnerID {
			if err = s.isFormAdmin(ctx, sub.ID, form.NamespaceID, payload.FormID); err != nil {
				return err
			}
		}
		_, err = s.forms.GetMember(ctx, payload.UserID, form.ID)
		if err == nil {
			return fun.ErrBadRequest("member is already a member of the form")
		}
		if !fun.Is(err, fun.CodeNotFound) {
			return err
		}
	}

	newMember := models.FormMember{
		UserID:  payload.UserID,
		FormID:  payload.FormID,
		Role:    payload.Role,
		AddedAt: time.Now(),
		AddedBy: sub.ID,
	}

	return s.forms.AddMember(ctx, newMember)
}
