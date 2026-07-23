package queries

import (
	"context"
	idx "sdk/identityx"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	ctx, span := q.tracer.Start(ctx, "ListMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := q.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.OwnerID != ident.Sub.ID {
		_, err := q.events.GetMember(ctx, event.ID, ident.Sub.ID)
		if err != nil {
			return nil, err
		}
	}

	members, err := q.events.ListEventMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
