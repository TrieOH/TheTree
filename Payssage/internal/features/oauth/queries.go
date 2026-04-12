package oauth

import (
	"context"
	"payssage/internal/platform/database"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/ports"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	workspaces   ports.WorkspaceRepo
	marketplaces ports.MarketplaceConfigRepo
	gaClient     *goauth.Client
	tx           database.TxRunner
	tracer       trace.Tracer
}

func NewQueryService(
	workspaces ports.WorkspaceRepo,
	marketplaces ports.MarketplaceConfigRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		workspaces:   workspaces,
		marketplaces: marketplaces,
		gaClient:     gaClient,
		tx:           tx,
		tracer:       tracer,
	}
}

func (uc *QueryService) ListMarketplaceConfigs(ctx context.Context, workspaceName string) ([]contracts.MarketplaceConfig, error) {
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
