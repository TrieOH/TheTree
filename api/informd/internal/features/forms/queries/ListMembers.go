package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *QueryService) ListMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := s.tracer.Start(ctx, "FormService.ListMembers")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	form, err := s.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if sub.ID != form.OwnerID {
		_, err := s.forms.GetMember(ctx, sub.ID, form.ID)
		if err != nil && fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	members, err := s.forms.ListDirectMembers(ctx, form.ID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
