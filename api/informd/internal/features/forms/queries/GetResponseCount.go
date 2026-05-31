package queries

import (
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *QueryService) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := s.tracer.Start(ctx, "FormService.GetResponseCount")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return 0, err
	}

	form, err := s.forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	if sub.ID != form.OwnerID {
		_, err := s.forms.GetMember(ctx, sub.ID, form.ID)
		if err != nil && fun.Is(err, fun.CodeNotFound) {
			return 0, err
		}
		if err != nil {
			return 0, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.forms.ResponsesCount(ctx, formID)
}
