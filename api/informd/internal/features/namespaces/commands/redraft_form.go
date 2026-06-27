package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Commands) ReDraftForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.ReDraftForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	form, err := s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, fun.ErrBadRequest("cannot redraft a form not on open")
	}

	count, err := s.forms.ResponsesCount(ctx, formID)
	if err != nil {
		return nil, err
	}

	if count != 0 {
		return nil, fun.ErrBadRequest("cannot redraft a form with responses")
	}

	if ident.Sub.ID != namespace.OwnerID {
		member, err := s.namespaces.GetMember(ctx, ident.Sub.ID, namespaceID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.NamespaceMemberRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.forms.ReDraft(ctx, form.ID)
}
