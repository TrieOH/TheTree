package queries

import (
	"context"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetFormResponseCount(ctx context.Context, namespaceID, formID uuid.UUID) (int, error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.GetFormResponseCount")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return 0, err
	}

	namespace, err := q.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return 0, err
	}

	if ident.Sub.ID != namespace.OwnerID {
		_, err = q.namespaces.GetMember(ctx, ident.Sub.ID, namespace.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return 0, err
		}
		if err != nil {
			_, err = q.forms.GetMember(ctx, ident.Sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return 0, err
			}
			if err != nil {
				return 0, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	return q.forms.ResponsesCount(ctx, formID)
}
