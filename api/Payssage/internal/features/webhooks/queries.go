package webhooks

import (
	"context"
	"payssage/internal/platform/database"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
	events     ports.WebhookEventRepo
	workspaces ports.WorkspaceRepo
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueryService(
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
	events ports.WebhookEventRepo,
	workspaces ports.WorkspaceRepo,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		endpoints:  endpoints,
		deliveries: deliveries,
		events:     events,
		workspaces: workspaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *QueryService) ListWebhookEndpoints(ctx context.Context, workspaceName string) ([]contracts.WebhookEndpoint, error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.ListWebhookEndpoints")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	return uc.endpoints.ListByWorkspace(ctx, workspace.ID)
}

func (uc *QueryService) ListWebhookDeliveries(ctx context.Context, workspaceName string, endpointID uuid.UUID) ([]contracts.WebhookDelivery, error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.ListWebhookDeliveries")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_webhook_deliveries"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.deliveries.ListByEndpoint(ctx, endpointID)
}

func (uc *QueryService) ListWebhookEvents(ctx context.Context, workspaceName string) ([]contracts.WebhookEventOriginal, error) {
	ctx, span := uc.tracer.Start(ctx, "QueryService.ListWebhookEvents")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_webhook_events"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.events.ListByWorkspace(ctx, workspace.ID)
}
