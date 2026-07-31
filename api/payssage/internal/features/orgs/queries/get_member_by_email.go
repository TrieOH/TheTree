package queries

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetMemberByEmail(ctx context.Context, email string, orgID uuid.UUID) (*idx.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "OrganizationService.GetMemberByEmail")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := q.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = q.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	actor, err := q.idx.Actors.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	_, err = q.orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	return actor, nil
}
