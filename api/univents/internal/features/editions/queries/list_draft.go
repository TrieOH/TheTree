package queries

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.ListDraft")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := q.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	member, err := q.events.GetMember(ctx, event.ID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return nil, err
	}
	if !member.Role.Minimum(models.EventMemberRoleStaff) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return q.editions.ListDraft(ctx, event.ID)
}
