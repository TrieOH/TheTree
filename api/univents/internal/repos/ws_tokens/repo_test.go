package ws_tokens_test

import (
	"context"
	"os"
	"testing"
	"time"

	"lib/testdb"
	"univents/internal/repos"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

// seedPurchase inserts the minimum purchase graph (edition → purchase) so a
// ws_token row has a valid FK target.
func seedPurchase(t *testing.T, q *sqlc.Queries) (purchaseID, userID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	event, err := q.CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:  uuid.New(),
		FullName: "Token Event",
		Slug:     "token-" + uuid.NewString()[:8],
		Status:   string(models.EventStatusActive),
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	edition, err := q.CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     event.ID,
		EditionName: "Token Edition",
		Slug:        "token-ed-" + uuid.NewString()[:8],
		StartsAt:    time.Now().Add(-time.Hour),
		EndsAt:      time.Now().Add(24 * time.Hour),
		CreatedBy:   event.OwnerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}

	userID = uuid.New()
	purchase, err := q.CreatePurchase(ctx, sqlc.CreatePurchaseParams{
		EditionID:   edition.ID,
		PurchaserID: userID,
		Status:      string(models.PurchaseStatusPending),
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		t.Fatalf("seed purchase: %v", err)
	}
	return purchase.ID, userID
}

func TestConsumeOneTime(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")

	r := repos.New(sqlc.New(pool))
	ctx := context.Background()
	purchaseID, userID := seedPurchase(t, sqlc.New(pool))

	hash := "hash-" + uuid.NewString()
	created, err := r.WsTokens.Create(ctx, &models.WsToken{
		PurchaseID: purchaseID,
		UserID:     userID,
		TokenHash:  hash,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.TokenHash != hash {
		t.Fatalf("hash = %q, want %q", created.TokenHash, hash)
	}
	if created.UsedAt != nil {
		t.Fatalf("fresh token already used")
	}

	consumed, err := r.WsTokens.Consume(ctx, hash)
	if err != nil {
		t.Fatalf("consume: %v", err)
	}
	if consumed == nil || consumed.ID != created.ID {
		t.Fatalf("consume = %v, want the created token", consumed)
	}

	// Second consume must miss the guard (already used) — one-time by construction.
	again, err := r.WsTokens.Consume(ctx, hash)
	if err != nil {
		t.Fatalf("second consume: %v", err)
	}
	if again != nil {
		t.Fatalf("second consume = %v, want nil (one-time)", again)
	}
}

func TestConsumeMissingOrExpired(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")

	r := repos.New(sqlc.New(pool))
	ctx := context.Background()
	purchaseID, userID := seedPurchase(t, sqlc.New(pool))

	// Missing hash → nil, not an error.
	missing, err := r.WsTokens.Consume(ctx, "no-such-hash")
	if err != nil {
		t.Fatalf("missing consume: %v", err)
	}
	if missing != nil {
		t.Fatalf("missing consume = %v, want nil", missing)
	}

	// Expired token → nil.
	hash := "expired-hash-" + uuid.NewString()
	_, err = r.WsTokens.Create(ctx, &models.WsToken{
		PurchaseID: purchaseID,
		UserID:     userID,
		TokenHash:  hash,
		ExpiresAt:  time.Now().Add(-time.Minute),
	})
	if err != nil {
		t.Fatalf("create expired: %v", err)
	}
	expired, err := r.WsTokens.Consume(ctx, hash)
	if err != nil {
		t.Fatalf("expired consume: %v", err)
	}
	if expired != nil {
		t.Fatalf("expired consume = %v, want nil", expired)
	}
}
