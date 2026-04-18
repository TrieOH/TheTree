package purchases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"univents/internal/platform/database"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	products  ports.ProductsRepository
	purchases ports.PurchaseRepository
	payments  *paymentsSDK.Client
	inventory ports.InventoryPublisher
	sessions  ports.PurchaseSessionStore
	ws        *sockets.Registry
	tracer    trace.Tracer
	tx        database.TxRunner
}

func NewAsynqService(
	products ports.ProductsRepository,
	purchases ports.PurchaseRepository,
	payments *paymentsSDK.Client,
	inventory ports.InventoryPublisher,
	sessions ports.PurchaseSessionStore,
	ws *sockets.Registry,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		products:  products,
		purchases: purchases,
		payments:  payments,
		inventory: inventory,
		sessions:  sessions,
		ws:        ws,
		tracer:    tracer,
		tx:        tx,
	}
}

func (uc *AsynqHandlers) HandleProductReservationExpiration(ctx context.Context, t *asynq.Task) error {
	var p contracts.ReservationExpiredPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	log.Printf("[task] reservation expired for session %s", p.SessionID)

	// 1. Check if purchase already exists for this session — if payment succeeded
	// before we got here, don't touch anything
	purchase, err := uc.purchases.GetBySessionID(ctx, p.SessionID)
	if err != nil && !errx.IsKind(err, "not_found") {
		return fmt.Errorf("failed to check purchase: %w", err)
	}
	if purchase != nil {
		log.Printf("[task] purchase already exists for session %s, skipping expiration", p.SessionID)
		return nil
	}

	// 2. Unreserve items — no-op if already unreserved
	var updates []contracts.InventoryUpdate
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var uErr error
		updates, uErr = uc.products.UnreserveItems(ctx, p.SessionID)
		if uErr != nil {
			return fmt.Errorf("failed to unreserve items: %w", uErr)
		}
		return nil
	}); err != nil {
		return err
	}

	if len(updates) > 0 {
		_ = uc.inventory.Publish(ctx, p.EditionID, updates)
	}

	// 3. Clean up Redis session if still there
	_ = uc.sessions.Delete(ctx, p.UserID, p.SessionID)

	// 4. Notify WS if still alive
	if err := uc.ws.Notify(p.SessionID.String(), sockets.WSMessage{
		Type:    "reservation_expired",
		Payload: "your reservation timed out",
	}); err != nil {
		log.Printf("[task] ws already closed for session %s: %v", p.SessionID, err)
	}

	uc.ws.Remove(p.SessionID.String())

	return nil
}
