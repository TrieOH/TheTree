package queries

import (
	"context"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := q.tracer.Start(ctx, "FormService.GetResponseCount")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return 0, err
	}

	form, err := q.forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	if ident.Sub.ID != form.OwnerID {
		_, err := q.forms.GetMember(ctx, ident.Sub.ID, form.ID)
		if err != nil && fun.Is(err, fun.CodeNotFound) {
			return 0, err
		}
		if err != nil {
			return 0, fun.ErrForbidden("insufficient permissions")
		}
	}

	return q.forms.ResponsesCount(ctx, formID)
}
