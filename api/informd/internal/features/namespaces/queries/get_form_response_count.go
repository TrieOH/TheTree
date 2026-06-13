package queries

import (
	"context"

	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetFormResponseCount(ctx context.Context, namespaceID, formID uuid.UUID) (int, error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.GetFormResponseCount")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return 0, err
	}

	namespace, err := q.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return 0, err
	}

	if sub.ID != namespace.OwnerID {
		_, err = q.namespaces.GetMember(ctx, sub.ID, namespace.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return 0, err
		}
		if err != nil {
			_, err = q.forms.GetMember(ctx, sub.ID, formID)
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
