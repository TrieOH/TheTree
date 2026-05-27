package commands

import (
	models2 "IdentityX/models"
	"Informd/models"
	"context"
	"time"

	"github.com/MintzyG/fun"
)

func (s *CommandService) AddFormMember(ctx context.Context, payload models.AddNamespaceFormMemberInput) error {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.AddFormMember")
	defer span.End()

	sub, err := models2.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to forms")
	}

	namespace, err := s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("owner of the namespace is already a member of the form")
	}

	if sub.ID != namespace.OwnerID {
		member, err := s.namespaces.GetMember(ctx, sub.ID, namespace.ID)
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
		return fun.ErrBadRequest("namespace member is already a member of the form")
	}

	form, err := s.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	_, err = s.forms.GetMember(ctx, payload.UserID, form.ID)
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
		AddedBy: sub.ID,
	}

	return s.forms.AddMember(ctx, newMember)
}
