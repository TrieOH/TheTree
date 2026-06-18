package webhooks

import (
	"encoding/json"
	"lib/telemetry"
	"log"
	"net/http"
	"time"

	"payssage/internal/shared/validation"

	_ "payssage/models"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

type RegisterWebhookEndpointRequest struct {
	URL string `json:"url" validate:"required"`
}

type WebhookEndpointResponse struct {
	ID          uuid.UUID  `json:"id"`
	ScopeID     uuid.UUID  `json:"scope_id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	URL         string     `json:"url"`
	Secret      string     `json:"secret"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// RegisterWebhookEndpoint godoc
// @Summary Register a webhook endpoint
// @Description Registers a URL to receive normalized payment events for the given workspace
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param request body RegisterWebhookEndpointRequest true "Endpoint details"
// @Success 201 {object} models.WebhookEndpoint "Endpoint registered successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/webhooks [post]
func (h *Handler) RegisterWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	var req RegisterWebhookEndpointRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	endpoint, err := h.commands.RegisterWebhookEndpoint(r.Context(), workspaceName, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(WebhookEndpointResponse{
		ID:          endpoint.ID,
		ScopeID:     endpoint.ScopeID,
		WorkspaceID: endpoint.WorkspaceID,
		URL:         endpoint.URL,
		Secret:      endpoint.Secret,
		CreatedAt:   endpoint.CreatedAt,
		DeletedAt:   endpoint.DeletedAt,
	}).Send(w)
}

// DeleteWebhookEndpoint godoc
// @Summary Delete a webhook endpoint
// @Description Deletes a registered webhook endpoint
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param endpoint_id path string true "Endpoint ID"
// @Success 200 {object} object "Endpoint deleted successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/webhooks/{endpoint_id} [delete]
func (h *Handler) DeleteWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	endpointID, rs := validation.GetUUID(r, "endpoint_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	if err := h.commands.DeleteWebhookEndpoint(r.Context(), workspaceName, endpointID); err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("endpoint deleted").Send(w)
}

type ProviderWebhookRequest struct {
	IntentID string `json:"intent_id" validate:"required"`
	Event    string `json:"event"     validate:"required"`
}

type MercadoPagoWebhookRequest struct {
	Action string `json:"action"`
	Data   struct {
		ID string `json:"id"`
	} `json:"data"`
}

// HandleProviderWebhook godoc
// @Summary Receive provider webhook
// @Description Receives a webhook from a payment provider, normalizes it and forwards to registered endpoints
// @Tags webhooks
// @Accept json
// @Produce json
// @Param provider path string true "Provider name (e.g. mock, stripe)"
// @Param request body ProviderWebhookRequest true "Provider webhook payload"
// @Success 200 {object} object "Received"
// @Failure 400 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /webhooks/{provider} [post]
func (h *Handler) HandleProviderWebhook(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	ctx := r.Context()
	var err error

	switch provider {
	case "mercadopago":
		if r.URL.Query().Get("type") != "payment" {
			telemetry.Log().Info("ignoring non-payment webhook", zap.String("type", r.URL.Query().Get("type")))
			fun.OK("ignored").Send(w)
			return
		}

		var req MercadoPagoWebhookRequest
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			fun.BadRequest("invalid payload").Send(w)
			return
		}

		log.Printf("[webhook] mercadopago received action=%s data.id=%s", req.Action, req.Data.ID)

		if req.Data.ID == "" {
			log.Printf("[webhook] mercadopago ignoring ping with no data.id")
			fun.OK("ignored").Send(w)
			return
		}

		fun.OK("received").Send(w)

		var rawPayload []byte
		rawPayload, err = json.Marshal(req)
		if err != nil {
			log.Printf("[webhook] mercadopago failed to marshal payload err=%v", err)
			rawPayload = []byte("{}")
		}

		var eventID uuid.UUID
		eventID, err = h.commands.CreateWebhookEvent(ctx, provider, req.Action, rawPayload)
		if err != nil {
			log.Printf("[webhook] failed to save event provider=%s err=%v", provider, err)
		}

		if err = h.commands.HandleMercadoPagoWebhook(ctx, req.Data.ID, eventID); err != nil {
			log.Printf("[webhook] mercadopago err=%v", err)
		}

	default:
		var req ProviderWebhookRequest
		if err = validation.ValidateInto(r, &req); err != nil {
			fun.Error(err).Send(w)
			return
		}

		fun.OK("received").Send(w)

		var rawPayload []byte
		rawPayload, err = json.Marshal(req)
		if err != nil {
			log.Printf("[webhook] mercadopago failed to marshal payload err=%v", err)
			rawPayload = []byte("{}")
		}

		var eventID uuid.UUID
		eventID, err = h.commands.CreateWebhookEvent(ctx, provider, req.Event, rawPayload)
		if err != nil {
			log.Printf("[webhook] failed to save event provider=%s err=%v", provider, err)
		}

		if err = h.commands.HandleProviderWebhook(ctx, eventID, provider, req.IntentID, req.Event); err != nil {
			log.Printf("[webhook] provider=%s err=%v", provider, err)
		}
	}
}

// ListWebhookEndpoints godoc
// @Summary List webhook endpoints
// @Description Lists all registered webhook endpoints for the given workspace
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} models.WebhookEndpoint "Endpoints retrieved successfully"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/webhooks [get]
func (h *Handler) ListWebhookEndpoints(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	endpoints, err := h.queries.ListWebhookEndpoints(r.Context(), workspaceName)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(endpoints).Send(w)
}

// ListWebhookDeliveries godoc
// @Summary List webhook deliveries
// @Description Lists all webhook deliveries for the given endpoint
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param endpoint_id path string true "Endpoint ID"
// @Success 200 {array} models.WebhookDelivery "Deliveries retrieved successfully"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/webhooks/{endpoint_id}/deliveries [get]
func (h *Handler) ListWebhookDeliveries(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	endpointID, rs := validation.GetUUID(r, "endpoint_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	deliveries, err := h.queries.ListWebhookDeliveries(r.Context(), workspaceName, endpointID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(deliveries).Send(w)
}

// ListWebhookEvents godoc
// @Summary List webhook events
// @Description Lists all inbound provider webhook events for the given workspace
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} models.WebhookEventOriginal "Events retrieved successfully"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/webhook-events [get]
func (h *Handler) ListWebhookEvents(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	events, err := h.queries.ListWebhookEvents(r.Context(), workspaceName)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(events).Send(w)
}
