package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) ListMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := q.tracer.Start(ctx, "FormService.ListMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	form, err := q.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != form.OwnerID {
		_, err := q.forms.GetMember(ctx, ident.Sub.ID, form.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	members, err := q.forms.ListDirectMembers(ctx, form.ID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
