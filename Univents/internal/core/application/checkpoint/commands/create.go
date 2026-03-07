package commands

import (
	"context"
	"encoding/json"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateCheckpointSpec) (out *domain.Checkpoint, err error) {
	ctx, span := uc.tracer.Start(ctx, "CheckpointService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var edition *domain.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	var validCheckpoint *domain.Checkpoint
	validCheckpoint, err = domain.NewCheckpoint(sub.ID, in, edition)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("checkpoints").
		Action("create").
		Scope(edition.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("checkpoint").SetMessage("insufficient permissions")
	}

	span.SetAttributes(attribute.String("checkpoint.id", validCheckpoint.ID.String()))

	meta := json.RawMessage(`{"color": "#fc620f", "icon": "FlagTriangleRight", "folder": "checkpoints"}`)
	var scope *goauth.Scope
	var idStr = validCheckpoint.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, validCheckpoint.Name, &idStr, &edition.GoauthScopeID, meta)
	if err != nil {
		return nil, err
	}
	validCheckpoint.AddScope(scope.ID)

	var created *domain.Checkpoint
	created, err = uc.checkpoints.Create(ctx, validCheckpoint)
	if err != nil {
		return nil, err
	}

	return created, nil
}
