package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []domain.Checkpoint, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "CheckpointService.AdminList")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	allowed, err := ga.Authz.Check().
		User(sub.ID).
		Object("checkpoints").
		Action("read").
		Scope(edition.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("checkpoints").SetMessage("insufficient permissions")
	}

	return uc.checkpoints.List(ctx, editionID)
}
