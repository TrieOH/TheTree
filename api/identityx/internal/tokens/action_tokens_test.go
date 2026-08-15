package tokens

import (
	"context"
	"errors"
	"testing"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

const testActionHMAC = "test-action-hmac"

// actionFixture wires an ActionTokenManager over a mockio ActionTokenRepo,
// capturing every record the manager persists, so tests can stub
// GetByJTI/Consume against the real shapes Mint produced.
type actionFixture struct {
	mgr      *ActionTokenManager
	repo     ports.ActionTokenRepo
	inserted []models.ActionToken
}

func newActionFixture(t *testing.T) *actionFixture {
	t.Helper()
	mock.SetUp(t)
	repo := mock.Mock[ports.ActionTokenRepo]()
	f := &actionFixture{repo: repo}
	_ = mock.When(repo.Insert(mock.AnyContext(), mock.Any[models.ActionToken]())).
		ThenAnswer(func(args []any) []any {
			rec := args[1].(models.ActionToken)
			f.inserted = append(f.inserted, rec)
			return []any{&rec, nil}
		})
	f.mgr = NewActionTokenManager(repo, []byte(testActionHMAC), ActionTokenConfig{
		VerifyTTL: 10 * time.Minute,
		ResetTTL:  5 * time.Minute,
	})
	return f
}

// stubLookup points GetByJTI at a specific record and Consume at a used
// variant of it (the shape Consume returns when it succeeds).
func (f *actionFixture) stubLookup(rec models.ActionToken) {
	used := rec
	now := time.Now()
	used.UsedAt = &now
	_ = mock.When(f.repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&rec, nil)
	_ = mock.When(f.repo.Consume(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)
}

func minted(t *testing.T, f *actionFixture, purpose models.ActionTokenPurpose, actorID uuid.UUID, projectID *uuid.UUID) (string, models.ActionToken) {
	t.Helper()
	token, _, err := f.mgr.Mint(context.Background(), purpose, actorID, projectID)
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	if len(f.inserted) == 0 {
		t.Fatalf("Mint did not persist an anti-replay record")
	}
	return token, f.inserted[len(f.inserted)-1]
}

func TestMintSignsAndPersists(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	projectID := uuid.New()

	token, rec := minted(t, f, models.EmailVerifyActionTokenPurpose, actorID, &projectID)

	claims := &models.ActionTokenClaims{}
	_, err := crypto.ParseHMACJWT(token, claims, []byte(testActionHMAC))
	if err != nil {
		t.Fatalf("token does not parse: %v", err)
	}
	if claims.Subject != actorID.String() {
		t.Fatalf("sub = %q, want %q", claims.Subject, actorID)
	}
	if claims.Purpose != string(models.EmailVerifyActionTokenPurpose) {
		t.Fatalf("purpose = %q", claims.Purpose)
	}
	if claims.ProjectID == nil || *claims.ProjectID != projectID {
		t.Fatalf("claim project_id not carried")
	}

	if rec.JTI.String() != claims.ID {
		t.Fatalf("persisted jti %q != claim id %q", rec.JTI, claims.ID)
	}
	if rec.Purpose != models.EmailVerifyActionTokenPurpose {
		t.Fatalf("persisted purpose = %q", rec.Purpose)
	}
	if rec.ActorID != actorID {
		t.Fatalf("persisted actor = %q", rec.ActorID)
	}
	if rec.UsedAt != nil {
		t.Fatalf("fresh record must be unused")
	}
}

func TestMintTTLPerPurpose(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()

	_, verifyTTL, err := f.mgr.Mint(context.Background(), models.EmailVerifyActionTokenPurpose, actorID, nil)
	if err != nil {
		t.Fatalf("Mint verify: %v", err)
	}
	if verifyTTL != 10*time.Minute {
		t.Fatalf("verify ttl = %v, want 10m", verifyTTL)
	}

	_, resetTTL, err := f.mgr.Mint(context.Background(), models.PasswordResetActionTokenPurpose, actorID, nil)
	if err != nil {
		t.Fatalf("Mint reset: %v", err)
	}
	if resetTTL != 5*time.Minute {
		t.Fatalf("reset ttl = %v, want 5m", resetTTL)
	}
}

func TestMintUnknownPurposeFails(t *testing.T) {
	f := newActionFixture(t)
	_, _, err := f.mgr.Mint(context.Background(), models.ActionTokenPurpose("bogus"), uuid.New(), nil)
	if err == nil || !fun.Is(err, fun.CodeInternal) {
		t.Fatalf("unknown purpose must fail internal, got %v", err)
	}
}

func TestRedeemSuccess(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	token, rec := minted(t, f, models.EmailVerifyActionTokenPurpose, actorID, nil)
	f.stubLookup(rec)

	got, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, token)
	if err != nil {
		t.Fatalf("Redeem: %v", err)
	}
	if got != actorID {
		t.Fatalf("redeemed actor = %q, want %q", got, actorID)
	}
	_, _ = mock.Verify(f.repo, mock.Times(1)).Consume(mock.AnyContext(), mock.Equal(rec.JTI))
}

func TestRedeemWrongPurpose(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	token, _ := minted(t, f, models.PasswordResetActionTokenPurpose, actorID, nil)

	_, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, token)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("wrong purpose must be bad request, got %v", err)
	}
	if !errors.Is(err, ErrActionTokenInvalid) {
		t.Fatalf("want ErrActionTokenInvalid, got %v", err)
	}
}

func TestRedeemGarbageToken(t *testing.T) {
	f := newActionFixture(t)
	_, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, "not-a-jwt")
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("garbage token must be bad request, got %v", err)
	}
	if !errors.Is(err, ErrActionTokenInvalid) {
		t.Fatalf("want ErrActionTokenInvalid, got %v", err)
	}
}

func TestRedeemUnknownRecord(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	token, _ := minted(t, f, models.EmailVerifyActionTokenPurpose, actorID, nil)

	_ = mock.When(f.repo.GetByJTI(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(nil, fun.ErrNotFound("no record"))

	_, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, token)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("unknown record must be bad request, got %v", err)
	}
	if !errors.Is(err, ErrActionTokenInvalid) {
		t.Fatalf("want ErrActionTokenInvalid, got %v", err)
	}
}

func TestRedeemUsed(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	token, rec := minted(t, f, models.EmailVerifyActionTokenPurpose, actorID, nil)

	used := rec
	now := time.Now()
	used.UsedAt = &now
	_ = mock.When(f.repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)

	got, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, token)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("used token must be bad request, got %v", err)
	}
	if !errors.Is(err, ErrActionTokenUsed) {
		t.Fatalf("want ErrActionTokenUsed, got %v", err)
	}
	if got != actorID {
		t.Fatalf("used error must still carry the actor id, got %q", got)
	}
	_, _ = mock.Verify(f.repo, mock.Times(0)).Consume(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestRedeemExpired(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	token, rec := minted(t, f, models.EmailVerifyActionTokenPurpose, actorID, nil)

	// The claim is still valid, but the anti-replay record has expired —
	// the record is the authority the token module checks.
	expired := rec
	expired.ExpiresAt = time.Now().Add(-time.Hour)
	_ = mock.When(f.repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&expired, nil)

	got, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, token)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("expired token must be bad request, got %v", err)
	}
	if !errors.Is(err, ErrActionTokenExpired) {
		t.Fatalf("want ErrActionTokenExpired, got %v", err)
	}
	if got != actorID {
		t.Fatalf("expired error must still carry the actor id, got %q", got)
	}
	_, _ = mock.Verify(f.repo, mock.Times(0)).Consume(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestRedeemConcurrentConsume(t *testing.T) {
	f := newActionFixture(t)
	actorID := uuid.New()
	token, rec := minted(t, f, models.EmailVerifyActionTokenPurpose, actorID, nil)

	_ = mock.When(f.repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&rec, nil)
	_ = mock.When(f.repo.Consume(mock.AnyContext(), mock.Equal(rec.JTI))).
		ThenReturn(nil, fun.ErrNotFound("consumed concurrently"))

	got, err := f.mgr.Redeem(context.Background(), models.EmailVerifyActionTokenPurpose, token)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("concurrent consume must be bad request, got %v", err)
	}
	if !errors.Is(err, ErrActionTokenUsed) {
		t.Fatalf("want ErrActionTokenUsed, got %v", err)
	}
	if got != actorID {
		t.Fatalf("concurrent consume must still carry the actor id, got %q", got)
	}
}

func TestDeleteExpiredDelegates(t *testing.T) {
	f := newActionFixture(t)
	_ = mock.When(f.repo.DeleteExpired(mock.AnyContext())).ThenReturn(nil)

	err := f.mgr.DeleteExpired(context.Background())
	if err != nil {
		t.Fatalf("DeleteExpired: %v", err)
	}
	_ = mock.Verify(f.repo, mock.Times(1)).DeleteExpired(mock.AnyContext())
}
