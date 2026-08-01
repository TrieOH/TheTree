// Package handlers implements the generated StrictServerInterface by
// delegating to one subpackage per feature. Auth, validation, and error
// mapping run in the strict middleware stack (see internal/app); the
// handlers here are pure domain logic + fun envelope construction.
package handlers

import (
	"context"

	"payssage/internal/handlers/collectors"
	"payssage/internal/handlers/intents"
	"payssage/internal/handlers/oauth"
	"payssage/internal/handlers/orgs"
	"payssage/internal/handlers/sellers"
	"payssage/internal/handlers/wallets"
	"payssage/internal/handlers/webhook_deliveries"
	"payssage/internal/handlers/webhook_endpoints"
	"payssage/internal/handlers/webhook_events"
	"payssage/internal/handlers/webhooks"
	"payssage/internal/openapi"
	"payssage/internal/services"
)

// Server implements openapi.StrictServerInterface.
type Server struct {
	orgs       *orgs.Handlers
	wallets    *wallets.Handlers
	oauth      *oauth.Handlers
	collectors *collectors.Handlers
	sellers    *sellers.Handlers
	intents    *intents.Handlers
	webhooks   *webhooks.Handlers
	endpoints  *webhook_endpoints.Handlers
	events     *webhook_events.Handlers
	deliveries *webhook_deliveries.Handlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		orgs:       orgs.New(ops.Organizations),
		wallets:    wallets.New(ops.Wallets),
		oauth:      oauth.New(ops.OAuth),
		collectors: collectors.New(ops.Collectors),
		sellers:    sellers.New(ops.Sellers),
		intents:    intents.New(ops.Intents),
		webhooks:   webhooks.New(ops.Webhooks),
		endpoints:  webhook_endpoints.New(ops.WebhookEndpoints),
		events:     webhook_events.New(ops.WebhookEvents),
		deliveries: webhook_deliveries.New(ops.WebhookDeliveries),
	}
}

// ── StrictServerInterface ────────────────────────────────────────────────

func (s *Server) ListOrganizations(ctx context.Context, req openapi.ListOrganizationsRequestObject) (openapi.ListOrganizationsResponseObject, error) {
	return s.orgs.ListOrganizations(ctx, req)
}

func (s *Server) CreateOrganization(ctx context.Context, req openapi.CreateOrganizationRequestObject) (openapi.CreateOrganizationResponseObject, error) {
	return s.orgs.CreateOrganization(ctx, req)
}

func (s *Server) ListOrganizationMembers(ctx context.Context, req openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	return s.orgs.ListOrganizationMembers(ctx, req)
}

func (s *Server) AddOrganizationMember(ctx context.Context, req openapi.AddOrganizationMemberRequestObject) (openapi.AddOrganizationMemberResponseObject, error) {
	return s.orgs.AddOrganizationMember(ctx, req)
}

func (s *Server) RemoveOrganizationMember(ctx context.Context, req openapi.RemoveOrganizationMemberRequestObject) (openapi.RemoveOrganizationMemberResponseObject, error) {
	return s.orgs.RemoveOrganizationMember(ctx, req)
}

func (s *Server) GetOrganizationMemberByID(ctx context.Context, req openapi.GetOrganizationMemberByIDRequestObject) (openapi.GetOrganizationMemberByIDResponseObject, error) {
	return s.orgs.GetOrganizationMemberByID(ctx, req)
}

func (s *Server) GetOrganizationMemberByEmail(ctx context.Context, req openapi.GetOrganizationMemberByEmailRequestObject) (openapi.GetOrganizationMemberByEmailResponseObject, error) {
	return s.orgs.GetOrganizationMemberByEmail(ctx, req)
}

func (s *Server) CreateWallet(ctx context.Context, req openapi.CreateWalletRequestObject) (openapi.CreateWalletResponseObject, error) {
	return s.wallets.CreateWallet(ctx, req)
}

func (s *Server) ListWallets(ctx context.Context, req openapi.ListWalletsRequestObject) (openapi.ListWalletsResponseObject, error) {
	return s.wallets.ListWallets(ctx, req)
}

func (s *Server) GetWallet(ctx context.Context, req openapi.GetWalletRequestObject) (openapi.GetWalletResponseObject, error) {
	return s.wallets.GetWallet(ctx, req)
}

func (s *Server) SetWalletFee(ctx context.Context, req openapi.SetWalletFeeRequestObject) (openapi.SetWalletFeeResponseObject, error) {
	return s.wallets.SetWalletFee(ctx, req)
}

func (s *Server) SetWalletSandbox(ctx context.Context, req openapi.SetWalletSandboxRequestObject) (openapi.SetWalletSandboxResponseObject, error) {
	return s.wallets.SetWalletSandbox(ctx, req)
}

func (s *Server) ListOrganizationWallets(ctx context.Context, req openapi.ListOrganizationWalletsRequestObject) (openapi.ListOrganizationWalletsResponseObject, error) {
	return s.wallets.ListOrganizationWallets(ctx, req)
}

func (s *Server) BindCollector(ctx context.Context, req openapi.BindCollectorRequestObject) (openapi.BindCollectorResponseObject, error) {
	return s.wallets.BindCollector(ctx, req)
}

func (s *Server) UnbindCollector(ctx context.Context, req openapi.UnbindCollectorRequestObject) (openapi.UnbindCollectorResponseObject, error) {
	return s.wallets.UnbindCollector(ctx, req)
}

func (s *Server) ListCollectors(ctx context.Context, req openapi.ListCollectorsRequestObject) (openapi.ListCollectorsResponseObject, error) {
	return s.collectors.ListCollectors(ctx, req)
}

func (s *Server) GetCollector(ctx context.Context, req openapi.GetCollectorRequestObject) (openapi.GetCollectorResponseObject, error) {
	return s.collectors.GetCollector(ctx, req)
}

func (s *Server) ListOrganizationCollectors(ctx context.Context, req openapi.ListOrganizationCollectorsRequestObject) (openapi.ListOrganizationCollectorsResponseObject, error) {
	return s.collectors.ListOrganizationCollectors(ctx, req)
}

func (s *Server) ListWalletSellers(ctx context.Context, req openapi.ListWalletSellersRequestObject) (openapi.ListWalletSellersResponseObject, error) {
	return s.sellers.ListWalletSellers(ctx, req)
}

func (s *Server) ListIntentsByProfile(ctx context.Context, req openapi.ListIntentsByProfileRequestObject) (openapi.ListIntentsByProfileResponseObject, error) {
	return s.intents.ListIntentsByProfile(ctx, req)
}

func (s *Server) GetIntent(ctx context.Context, req openapi.GetIntentRequestObject) (openapi.GetIntentResponseObject, error) {
	return s.intents.GetIntent(ctx, req)
}

func (s *Server) CancelIntent(ctx context.Context, req openapi.CancelIntentRequestObject) (openapi.CancelIntentResponseObject, error) {
	return s.intents.CancelIntent(ctx, req)
}

func (s *Server) ListWalletIntents(ctx context.Context, req openapi.ListWalletIntentsRequestObject) (openapi.ListWalletIntentsResponseObject, error) {
	return s.intents.ListWalletIntents(ctx, req)
}

func (s *Server) ListOrganizationIntents(ctx context.Context, req openapi.ListOrganizationIntentsRequestObject) (openapi.ListOrganizationIntentsResponseObject, error) {
	return s.intents.ListOrganizationIntents(ctx, req)
}

func (s *Server) Checkout(ctx context.Context, req openapi.CheckoutRequestObject) (openapi.CheckoutResponseObject, error) {
	return s.intents.Checkout(ctx, req)
}

func (s *Server) HardCreateIntent(ctx context.Context, req openapi.HardCreateIntentRequestObject) (openapi.HardCreateIntentResponseObject, error) {
	return s.intents.HardCreateIntent(ctx, req)
}

func (s *Server) ConnectProvider(ctx context.Context, req openapi.ConnectProviderRequestObject) (openapi.ConnectProviderResponseObject, error) {
	return s.oauth.ConnectProvider(ctx, req)
}

func (s *Server) ProviderCallback(ctx context.Context, req openapi.ProviderCallbackRequestObject) (openapi.ProviderCallbackResponseObject, error) {
	return s.oauth.ProviderCallback(ctx, req)
}

func (s *Server) RevokeProvider(ctx context.Context, req openapi.RevokeProviderRequestObject) (openapi.RevokeProviderResponseObject, error) {
	return s.oauth.RevokeProvider(ctx, req)
}

func (s *Server) ReceiveWebhook(ctx context.Context, req openapi.ReceiveWebhookRequestObject) (openapi.ReceiveWebhookResponseObject, error) {
	return s.webhooks.ReceiveWebhook(ctx, req)
}

func (s *Server) CreateWebhookEndpoint(ctx context.Context, req openapi.CreateWebhookEndpointRequestObject) (openapi.CreateWebhookEndpointResponseObject, error) {
	return s.endpoints.CreateWebhookEndpoint(ctx, req)
}

func (s *Server) ListWebhookEndpoints(ctx context.Context, req openapi.ListWebhookEndpointsRequestObject) (openapi.ListWebhookEndpointsResponseObject, error) {
	return s.endpoints.ListWebhookEndpoints(ctx, req)
}

func (s *Server) GetWebhookEndpoint(ctx context.Context, req openapi.GetWebhookEndpointRequestObject) (openapi.GetWebhookEndpointResponseObject, error) {
	return s.endpoints.GetWebhookEndpoint(ctx, req)
}

func (s *Server) DeleteWebhookEndpoint(ctx context.Context, req openapi.DeleteWebhookEndpointRequestObject) (openapi.DeleteWebhookEndpointResponseObject, error) {
	return s.endpoints.DeleteWebhookEndpoint(ctx, req)
}

func (s *Server) ListWebhookEvents(ctx context.Context, req openapi.ListWebhookEventsRequestObject) (openapi.ListWebhookEventsResponseObject, error) {
	return s.events.ListWebhookEvents(ctx, req)
}

func (s *Server) GetWebhookEvent(ctx context.Context, req openapi.GetWebhookEventRequestObject) (openapi.GetWebhookEventResponseObject, error) {
	return s.events.GetWebhookEvent(ctx, req)
}

func (s *Server) ListWebhookDeliveries(ctx context.Context, req openapi.ListWebhookDeliveriesRequestObject) (openapi.ListWebhookDeliveriesResponseObject, error) {
	return s.deliveries.ListWebhookDeliveries(ctx, req)
}

func (s *Server) GetWebhookDelivery(ctx context.Context, req openapi.GetWebhookDeliveryRequestObject) (openapi.GetWebhookDeliveryResponseObject, error) {
	return s.deliveries.GetWebhookDelivery(ctx, req)
}
