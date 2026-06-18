package intents

import (
	"context"
	"lib/authz"
	"payssage/models"
	"payssage/ports"

	"lib/database"
	"payssage/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	intents    ports.IntentRepository
	workspaces ports.WorkspaceRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueryService(
	intents ports.IntentRepository,
	workspaces ports.WorkspaceRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		intents:    intents,
		workspaces: workspaces,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *QueryService) GetByID(ctx context.Context, id uuid.UUID) (intent *models.Intent, err error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.GetByID")
	defer span.End()

	ws, err := authz.RequireWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	intent, err = uc.intents.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if intent.WorkspaceID != ws.ID {
		return nil, errx.Forbidden("intent").SetMessage("intent does not belong to this workspace")
	}

	return intent, nil
}

func (uc *QueryService) List(ctx context.Context) (intents []models.Intent, err error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.List")
	defer span.End()

	// try workspace from API key first
	ws, err := authz.RequireWorkspace(ctx)
	if err == nil {
		return uc.intents.ListIntentsByWorkspace(ctx, ws.ID)
	}

	// fall back to user session — list all workspaces then all intents
	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspaces, err := uc.workspaces.List(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	for _, w := range workspaces {
		wsIntents, err := uc.intents.ListIntentsByWorkspace(ctx, w.ID)
		if err != nil {
			return nil, err
		}
		intents = append(intents, wsIntents...)
	}

	return intents, nil
}

func (uc *QueryService) ListByWorkspace(ctx context.Context, wsName string) (intents []models.Intent, err error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.List")
	defer span.End()

	// try workspace from API key first
	ws, err := authz.RequireWorkspace(ctx)
	if err == nil {
		return uc.intents.ListIntentsByWorkspace(ctx, ws.ID)
	}

	// fall back to user session — list all workspaces then all intents
	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, wsName, sub.ID)
	if err != nil {
		return nil, err
	}

	intents, err = uc.intents.ListIntentsByWorkspace(ctx, workspace.ID)
	if err != nil {
		return nil, err
	}

	return intents, nil
}
