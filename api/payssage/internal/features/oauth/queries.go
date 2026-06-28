package oauth

import (
	"context"
	"payssage/ports"

	"lib/authz"
	"lib/database"
	"payssage/models"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	workspaces   ports.WorkspaceRepo
	marketplaces ports.MarketplaceConfigRepo
	logger       *zap.Logger
	tx           database.TxRunner
	tracer       trace.Tracer
}

func NewQueryService(
	workspaces ports.WorkspaceRepo,
	marketplaces ports.MarketplaceConfigRepo,
	logger *zap.Logger,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		workspaces:   workspaces,
		marketplaces: marketplaces,
		logger:       logger,
		tx:           tx,
		tracer:       tracer,
	}
}

func (uc *QueryService) ListMarketplaceConfigs(ctx context.Context, workspaceName string) ([]models.MarketplaceConfig, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.ListMarketplaceConfigs")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	return uc.marketplaces.List(ctx, workspace.ID)
}
