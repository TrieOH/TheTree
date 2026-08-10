package payments_test

import (
	"context"
	"os"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	idx "sdk/identityx"
	payssage "sdk/payssage"

	"univents/internal/authz"
	"univents/internal/services/payments"
	"univents/models"
	"univents/ports"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

func adminCtx() context.Context {
	return idx.WithIdentity(context.Background(), &idx.Identity{Sub: idx.Subject{ID: uuid.New()}})
}

// platformWalletID is the shared wallet every test OAuths onto (D6).
func platformWalletID() uuid.UUID {
	return uuid.MustParse("11111111-1111-1111-1111-111111111111")
}

func newOps(repo ports.EventRepo, pc payments.PayssageClient) *payments.Operations {
	return payments.NewOperations(repo, pc, authz.New(repo), platformWalletID())
}

func stubAdminEvent(repo ports.EventRepo, event *models.Event) {
	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleAdmin, nil)
}

func TestConnect_OAuthsOntoPlatformWallet(t *testing.T) {
	mock.SetUp(t)
	repo := mock.Mock[ports.EventRepo]()
	pc := mock.Mock[payments.PayssageClient]()

	eventID := uuid.New()
	stubAdminEvent(repo, &models.Event{ID: eventID, Slug: "my-event"})

	mock.When(pc.ConnectProvider(mock.AnyContext(), mock.Any[string](), mock.Any[payssage.ConnectProviderRequest]())).
		ThenReturn("https://auth.mercadopago.com/consent", nil)

	t.Setenv("APP_URL", "https://events.test")

	res, err := newOps(repo, pc).Connect(adminCtx(), eventID, "mercado_pago")
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if res.AuthURL != "https://auth.mercadopago.com/consent" {
		t.Fatalf("auth_url = %q", res.AuthURL)
	}
	if res.WalletID != platformWalletID() {
		t.Fatalf("wallet_id = %s, want platform wallet %s", res.WalletID, platformWalletID())
	}

	// No wallet is created or configured — the platform wallet is
	// env-configured (D6); the seller OAuths straight onto it.
	_, _ = mock.Verify(repo, mock.Never()).SetPaymentsConfig(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[*uuid.UUID](), mock.Any[*string]())

	// Callback URL is per event: event id in the path, no query — so
	// payssage's own `?credential_id=…&public_key=…` append works untouched.
	_, _ = mock.Verify(pc, mock.Once()).ConnectProvider(mock.AnyContext(), mock.Exact("mercado_pago"),
		mock.Equal(payssage.ConnectProviderRequest{
			Flow:             payssage.OAuthFlowSeller,
			WalletID:         func() *uuid.UUID { w := platformWalletID(); return &w }(),
			FinalRedirectURL: "https://events.test/events/" + eventID.String() + "/payssage/oauth/callback",
		}))
}

func TestConnect_ForbiddenForNonAdmin(t *testing.T) {
	mock.SetUp(t)
	repo := mock.Mock[ports.EventRepo]()
	pc := mock.Mock[payments.PayssageClient]()

	eventID := uuid.New()
	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(&models.Event{ID: eventID}, nil)
	// staff rank (0) is below admin (1) → forbidden
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleStaff, nil)

	_, err := newOps(repo, pc).Connect(adminCtx(), eventID, "mercado_pago")
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
	_, _ = mock.Verify(pc, mock.Never()).ConnectProvider(mock.AnyContext(), mock.Any[string](), mock.Any[payssage.ConnectProviderRequest]())
}

func TestComplete_PersistsSellerAndPublicKey(t *testing.T) {
	mock.SetUp(t)
	repo := mock.Mock[ports.EventRepo]()
	pc := mock.Mock[payments.PayssageClient]()

	eventID := uuid.New()
	sellerID := uuid.New()
	stubAdminEvent(repo, &models.Event{ID: eventID})

	mock.When(pc.ListWalletSellers(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]payssage.Seller{{ID: sellerID, Provider: "mercado_pago"}}, nil)
	var gotSeller *uuid.UUID
	var gotKey *string
	mock.When(repo.SetPaymentsConfig(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[*uuid.UUID](), mock.Any[*string]())).
		ThenAnswer(func(args []any) []any {
			gotSeller = args[2].(*uuid.UUID)
			gotKey = args[3].(*string)
			return []any{&models.Event{ID: eventID, PayssageSellerID: &sellerID}, nil}
		})

	event, err := newOps(repo, pc).Complete(adminCtx(), eventID, sellerID, "TEST-PUBLIC-KEY")
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if event.PayssageSellerID == nil || *event.PayssageSellerID != sellerID {
		t.Fatalf("seller not persisted: %+v", event)
	}
	if gotSeller == nil || *gotSeller != sellerID {
		t.Fatalf("persisted seller = %v, want %s", gotSeller, sellerID)
	}
	if gotKey == nil || *gotKey != "TEST-PUBLIC-KEY" {
		t.Fatalf("persisted public key = %v", gotKey)
	}

	// The seller is verified against the platform wallet (D6).
	_, _ = mock.Verify(pc, mock.Once()).ListWalletSellers(mock.AnyContext(), mock.Equal(platformWalletID()))
}

func TestComplete_RejectsSellerNotInWallet(t *testing.T) {
	mock.SetUp(t)
	repo := mock.Mock[ports.EventRepo]()
	pc := mock.Mock[payments.PayssageClient]()

	eventID := uuid.New()
	foreignSeller := uuid.New()
	stubAdminEvent(repo, &models.Event{ID: eventID})

	mock.When(pc.ListWalletSellers(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]payssage.Seller{{ID: uuid.New(), Provider: "mercado_pago"}}, nil)

	_, err := newOps(repo, pc).Complete(adminCtx(), eventID, foreignSeller, "k")
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
	_, _ = mock.Verify(repo, mock.Never()).SetPaymentsConfig(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[*uuid.UUID](), mock.Any[*string]())
}

func TestDisconnect_UnlinksSellerAndKey(t *testing.T) {
	mock.SetUp(t)
	repo := mock.Mock[ports.EventRepo]()
	pc := mock.Mock[payments.PayssageClient]()

	eventID := uuid.New()
	stubAdminEvent(repo, &models.Event{ID: eventID})

	mock.When(repo.ClearPaymentsConfig(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Event{ID: eventID}, nil)

	event, err := newOps(repo, pc).Disconnect(adminCtx(), eventID)
	if err != nil {
		t.Fatalf("Disconnect: %v", err)
	}
	if event.ID != eventID {
		t.Fatalf("event id = %s, want %s", event.ID, eventID)
	}
	_, _ = mock.Verify(repo, mock.Once()).ClearPaymentsConfig(mock.AnyContext(), mock.Equal(eventID))
}
