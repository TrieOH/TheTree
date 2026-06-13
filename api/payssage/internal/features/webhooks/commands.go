package webhooks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"payssage/internal/platform/database"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	endpoints   ports.WebhookEndpointRepo
	deliveries  ports.WebhookDeliveryRepo
	events      ports.WebhookEventRepo
	workspaces  ports.WorkspaceRepo
	intents     ports.IntentRepository
	credentials ports.ProviderCredentialRepo
	asynq       *asynq.Client
	az          *authzed.Client
	tx          database.TxRunner
	tracer      trace.Tracer
}

func NewCommandService(
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
	events ports.WebhookEventRepo,
	workspaces ports.WorkspaceRepo,
	intents ports.IntentRepository,
	credentials ports.ProviderCredentialRepo,
	asynq *asynq.Client,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		endpoints:   endpoints,
		deliveries:  deliveries,
		events:      events,
		workspaces:  workspaces,
		intents:     intents,
		credentials: credentials,
		asynq:       asynq,
		az:          az,
		tx:          tx,
		tracer:      tracer,
	}
}

func (uc *CommandService) RegisterWebhookEndpoint(ctx context.Context, workspaceName, url string) (*contracts.WebhookEndpoint, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RegisterWebhookEndpoint")
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
		authz.Permission("create_webhooks"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return nil, err
	}

	// generate HMAC secret
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, err
	}
	secret := hex.EncodeToString(secretBytes)

	endpoint, err := contracts.NewWebhookEndpoint(workspace.ID, url, secret)
	if err != nil {
		return nil, err
	}

	created, err := uc.endpoints.Create(ctx, *endpoint)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) CreateWebhookEvent(ctx context.Context, provider, eventType string, payload json.RawMessage) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	event := contracts.WebhookEventOriginal{
		ID:        id,
		Provider:  provider,
		EventType: eventType,
		Payload:   payload,
	}
	created, err := uc.events.Create(ctx, event)
	if err != nil {
		return uuid.Nil, err
	}
	return created.ID, nil
}

func (uc *CommandService) DeleteWebhookEndpoint(ctx context.Context, workspaceName string, endpointID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.DeleteWebhookEndpoint")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("delete_webhooks"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return err
	}

	return uc.endpoints.Delete(ctx, endpointID, workspace.ID)
}

func (uc *CommandService) Dispatch(ctx context.Context, provider, intentID, event string, eventID uuid.UUID) error {
	return uc.HandleProviderWebhook(ctx, eventID, provider, intentID, event)
}

func (uc *CommandService) EnrichWebhookEvent(ctx context.Context, eventID, workspaceID, intentID uuid.UUID, externalID string) {
	if _, err := uc.events.Enrich(ctx, eventID, workspaceID, intentID, externalID); err != nil {
		log.Printf("[webhook_event] failed to enrich event=%s err=%v", eventID, err)
	}
}

func (uc *CommandService) HandleMercadoPagoWebhook(ctx context.Context, mpOrderID string, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.HandleMercadoPagoWebhook")
	defer span.End()

	log.Printf("[mp-webhook] looking up intent for order_id=%s", mpOrderID)

	intent, err := uc.intents.GetByMPOrderID(ctx, mpOrderID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			log.Printf("[mp-webhook] no intent found for mp order %s", mpOrderID)
			return nil
		}
		return err
	}

	log.Printf("[mp-webhook] found intent=%s status=%s", intent.ID, intent.Status)

	if intent.SellerCredentialID == nil {
		return fmt.Errorf("intent %s has no seller credential", intent.ID)
	}

	cred, err := uc.credentials.GetByID(ctx, *intent.SellerCredentialID)
	if err != nil {
		return fmt.Errorf("failed to fetch seller credential: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.mercadopago.com/v1/payments/"+mpOrderID,
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cred.Credentials.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[mp-webhook] failed to fetch mp order %s: %v", mpOrderID, err)
		return err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("mercadopago fetch order error %d: %s", resp.StatusCode, string(rawBody))
	}

	var mpOrder struct {
		ID           int64  `json:"id"`
		Status       string `json:"status"`
		StatusDetail string `json:"status_detail"`
	}
	if err := json.Unmarshal(rawBody, &mpOrder); err != nil {
		return err
	}

	log.Printf("[mp-webhook] mp order %s status=%s status_detail=%s", mpOrderID, mpOrder.Status, mpOrder.StatusDetail)

	event := mapMPOrderStatusToEvent(mpOrder.Status, mpOrder.StatusDetail)
	if event == "" {
		log.Printf("[mp-webhook] ignoring mp order status=%s detail=%s order=%s", mpOrder.Status, mpOrder.StatusDetail, mpOrderID)
		return nil
	}

	log.Printf("[mp-webhook] dispatching event=%s for intent=%s", event, intent.ID)
	return uc.HandleProviderWebhook(ctx, eventID, "mercadopago", intent.ID.String(), event)
}

func mapMPOrderStatusToEvent(status, statusDetail string) string {
	switch status {
	case "approved":
		return contracts.EventPaymentSucceeded
	case "pending", "in_process", "authorized":
		return ""
	case "rejected", "cancelled", "refunded", "charged_back":
		return contracts.EventPaymentFailed
	default:
		return ""
	}
}

func (uc *CommandService) HandleProviderWebhook(ctx context.Context, eventID uuid.UUID, provider, intentID string, event string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.HandleProviderWebhook")
	defer span.End()

	var id uuid.UUID
	id, err = uuid.Parse(intentID)
	if err != nil {
		return errx.Invalid("intent").SetMessage("invalid intent_id")
	}

	var intent *contracts.Intent
	intent, err = uc.intents.GetByID(ctx, id)
	if err != nil {
		return err
	}

	switch event {
	case contracts.EventPaymentSucceeded:
		if !alreadyInTargetState(event, intent.Status) {
			log.Printf("[webhook] confirming intent=%s", id)
			intent, err = uc.intents.Confirm(ctx, id)
		} else {
			log.Printf("[webhook] intent=%s already in target state, skipping mutation", id)
		}
	case contracts.EventPaymentFailed:
		if !alreadyInTargetState(event, intent.Status) {
			log.Printf("[webhook] failing intent=%s", id)
			intent, err = uc.intents.Fail(ctx, id)
		} else {
			log.Printf("[webhook] intent=%s already in target state, skipping mutation", id)
		}
	case contracts.EventPaymentCancelled:
		if !alreadyInTargetState(event, intent.Status) {
			log.Printf("[webhook] cancelling intent=%s", id)
			intent, err = uc.intents.Cancel(ctx, id)
		} else {
			log.Printf("[webhook] intent=%s already in target state, skipping mutation", id)
		}
	default:
		return errx.Invalid("event").SetMessage("unknown event type: " + event)
	}
	if err != nil {
		log.Printf("[webhook] failed to update intent=%s event=%s err=%v", id, event, err)
		return err
	}

	log.Printf("[webhook] intent=%s status=%s", intent.ID, intent.Status)

	// build normalized payload
	var payloadBytes []byte
	payloadBytes, err = json.Marshal(contracts.WebhookPayload{
		Event:           event,
		IntentID:        intent.ID,
		WorkspaceID:     intent.WorkspaceID,
		Amount:          intent.Amount,
		Currency:        intent.Currency,
		Provider:        intent.Provider,
		Metadata:        intent.Metadata,
		MercadoPagoData: intent.MercadoPagoData,
	})
	if err != nil {
		return err
	}

	// fetch all registered endpoints for this workspace
	var endpoints []contracts.WebhookEndpoint
	endpoints, err = uc.endpoints.ListByWorkspace(ctx, intent.WorkspaceID)
	if err != nil {
		return err
	}

	log.Printf("[webhook] found %d endpoints for workspace=%s", len(endpoints), intent.WorkspaceID)

	// enqueue delivery task per endpoint
	for _, endpoint := range endpoints {
		var delivery *contracts.WebhookDelivery
		delivery, err = contracts.NewWebhookDelivery(endpoint.ID, intent.ID, event, payloadBytes)
		if err != nil {
			log.Printf("[webhook] failed to create delivery object for endpoint %s: %v", endpoint.ID, err)
			continue
		}
		var created *contracts.WebhookDelivery
		created, err = uc.deliveries.Create(ctx, *delivery)
		if err != nil {
			log.Printf("[webhook] failed to create delivery record for endpoint %s: %v", endpoint.ID, err)
			continue
		}

		var task *asynq.Task
		task, err = contracts.NewDeliverWebhookTask(created.ID, endpoint.ID, endpoint.URL, endpoint.Secret, payloadBytes)
		if err != nil {
			log.Printf("[webhook] failed to create delivery task for endpoint %s: %v", endpoint.ID, err)
			continue
		}

		if _, err = uc.asynq.EnqueueContext(context.Background(), task); err != nil {
			log.Printf("[webhook] failed to enqueue delivery task for endpoint %s: %v", endpoint.ID, err)
		} else {
			log.Printf("[webhook] enqueued delivery for endpoint=%s url=%s", endpoint.ID, endpoint.URL)
		}
	}

	if eventID != uuid.Nil {
		uc.EnrichWebhookEvent(ctx, eventID, intent.WorkspaceID, intent.ID, intentID)
	}

	return nil
}

func alreadyInTargetState(event string, status contracts.IntentStatus) bool {
	switch event {
	case contracts.EventPaymentSucceeded:
		return status == contracts.IntentStatusSucceeded
	case contracts.EventPaymentFailed:
		return status == contracts.IntentStatusFailed
	case contracts.EventPaymentCancelled:
		return status == contracts.IntentStatusCancelled
	}
	return false
}
