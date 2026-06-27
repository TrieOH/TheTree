package commands

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (s *Commands) AddFormMember(ctx context.Context, payload models.AddNamespaceFormMemberInput) error {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.AddFormMember")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return fun.ErrBadRequest("users can't add themselves to forms")
	}

	namespace, err := s.namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return fun.ErrBadRequest("owner of the namespace is already a member of the form")
	}

	if ident.Sub.ID != namespace.OwnerID {
		member, err := s.namespaces.GetMember(ctx, ident.Sub.ID, namespace.ID)
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
		AddedBy: ident.Sub.ID,
	}

	return s.forms.AddMember(ctx, newMember)
}
