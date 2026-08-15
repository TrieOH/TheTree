package authn

import (
	"context"
	"testing"
	"time"

	"IdentityX/internal/tokens"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

const testHMAC = "test-hmac-secret"

// emailOps builds authn operations with the real token managers over
// per-test mockio repos, exactly like production wiring.
func emailOps(t *testing.T, actors ports.ActorRepo, actionTokens *tokens.ActionTokenManager) *Operations {
	t.Helper()
	return NewOperations(
		actors,
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.PlatformRolesRepo](),
		tokens.NewManager(mock.Mock[ports.CryptoKeysRepo](), mock.Mock[ports.BlacklistRepo](), mock.Mock[ports.ActorRepo](), mock.Mock[ports.ProjectRepo](), tokens.Config{}),
		actionTokens,
		mock.Mock[ports.EmailSender](),
	)
}

// mintFor mints through the ActionTokenManager with the repo's Insert
// stubbed, and returns the token plus the anti-replay record the manager
// persisted — the exact shape tests stub GetByJTI/Consume against.
func mintFor(t *testing.T, mgr *tokens.ActionTokenManager, repo ports.ActionTokenRepo, actorID uuid.UUID, purpose models.ActionTokenPurpose) (string, models.ActionToken) {
	t.Helper()
	var rec models.ActionToken
	_ = mock.When(repo.Insert(mock.AnyContext(), mock.Any[models.ActionToken]())).
		ThenAnswer(func(args []any) []any {
			rec = args[1].(models.ActionToken)
			return []any{&rec, nil}
		})
	token, _, err := mgr.Mint(context.Background(), purpose, actorID, nil)
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	return token, rec
}

// actionMgr builds an ActionTokenManager over the repo with the standard
// test TTLs.
func actionMgr(repo ports.ActionTokenRepo) *tokens.ActionTokenManager {
	return tokens.NewActionTokenManager(repo, []byte(testHMAC), tokens.ActionTokenConfig{
		VerifyTTL: 10 * time.Minute,
		ResetTTL:  10 * time.Minute,
	})
}

func TestVerifyEmailSuccess(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	actorID := actor.ID

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, rec := mintFor(t, mgr, repo, actorID, models.EmailVerifyActionTokenPurpose)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.SetVerifiedAt(mock.AnyContext(), mock.Equal(actorID), mock.Any[time.Time]())).ThenReturn(nil)

	_ = mock.When(repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&rec, nil)
	used := rec
	now := time.Now()
	used.UsedAt = &now
	_ = mock.When(repo.Consume(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)

	ops := emailOps(t, actors, mgr)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	_, _ = mock.Verify(repo, mock.Times(1)).Consume(mock.AnyContext(), mock.Equal(rec.JTI))
	_ = mock.Verify(actors, mock.Times(1)).SetVerifiedAt(mock.AnyContext(), mock.Equal(actorID), mock.Any[time.Time]())
}

func TestVerifyEmailIdempotentWhenAlreadyVerified(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	now := time.Now()
	actor.VerifiedAt = &now

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, rec := mintFor(t, mgr, repo, actor.ID, models.EmailVerifyActionTokenPurpose)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)

	used := rec
	used.UsedAt = &now
	_ = mock.When(repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)

	ops := emailOps(t, actors, mgr)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err != nil {
		t.Fatalf("re-click on verified account must succeed, got %v", err)
	}
}

func TestVerifyEmailRejectsUsedTokenOnUnverifiedAccount(t *testing.T) {
	mock.SetUp(t)
	actor := testActor() // unverified

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, rec := mintFor(t, mgr, repo, actor.ID, models.EmailVerifyActionTokenPurpose)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)

	now := time.Now()
	used := rec
	used.UsedAt = &now
	_ = mock.When(repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)

	ops := emailOps(t, actors, mgr)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for replayed token, got %v", err)
	}
}

func TestVerifyEmailRejectsWrongPurpose(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, _ := mintFor(t, mgr, repo, actor.ID, models.PasswordResetActionTokenPurpose)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), mgr)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for reset token on verify, got %v", err)
	}
}

func TestVerifyEmailRejectsGarbageToken(t *testing.T) {
	mock.SetUp(t)
	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), mgr)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: "not-a-jwt"})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for garbage token, got %v", err)
	}
}

func TestResetPasswordSuccess(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, rec := mintFor(t, mgr, repo, actor.ID, models.PasswordResetActionTokenPurpose)

	actors := mock.Mock[ports.ActorRepo]()
	captor := mock.Captor[string]()
	_ = mock.When(actors.UpdatePasswordHash(mock.AnyContext(), mock.Equal(actor.ID), captor.Capture())).ThenReturn(nil)

	_ = mock.When(repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&rec, nil)
	used := rec
	now := time.Now()
	used.UsedAt = &now
	_ = mock.When(repo.Consume(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)

	ops := emailOps(t, actors, mgr)
	newPassword := "NewPassword123!"
	err := ops.ResetPassword(context.Background(), models.ResetPasswordInput{Token: token, NewPassword: newPassword})
	if err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	hashes := captor.Values()
	if len(hashes) != 1 {
		t.Fatalf("UpdatePasswordHash called %d times, want 1", len(hashes))
	}
	verifyErr := crypto.Verify(newPassword, hashes[0])
	if verifyErr != nil {
		t.Fatalf("stored hash does not verify against new password: %v", verifyErr)
	}
	_, _ = mock.Verify(repo, mock.Times(1)).Consume(mock.AnyContext(), mock.Equal(rec.JTI))
}

func TestResetPasswordRejectsUsedToken(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, rec := mintFor(t, mgr, repo, actor.ID, models.PasswordResetActionTokenPurpose)

	now := time.Now()
	used := rec
	used.UsedAt = &now
	_ = mock.When(repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&used, nil)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), mgr)
	err := ops.ResetPassword(context.Background(), models.ResetPasswordInput{Token: token, NewPassword: "NewPassword123!"})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for used reset token, got %v", err)
	}
}

func TestResetPasswordRejectsExpiredToken(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()

	repo := mock.Mock[ports.ActionTokenRepo]()
	mgr := actionMgr(repo)
	token, rec := mintFor(t, mgr, repo, actor.ID, models.PasswordResetActionTokenPurpose)

	// The claim is still valid, but the anti-replay record has expired —
	// the record is the authority the token module checks.
	expired := rec
	expired.ExpiresAt = time.Now().Add(-time.Hour)
	_ = mock.When(repo.GetByJTI(mock.AnyContext(), mock.Equal(rec.JTI))).ThenReturn(&expired, nil)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), mgr)
	err := ops.ResetPassword(context.Background(), models.ResetPasswordInput{Token: token, NewPassword: "NewPassword123!"})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for expired reset token, got %v", err)
	}
}

func TestResendVerificationSilentNoOps(t *testing.T) {
	mock.SetUp(t)
	email := "user@example.com"

	// unknown email → 200, nothing dispatched
	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByEmail(mock.AnyContext(), mock.Equal(email), mock.Any[*uuid.UUID]())).
		ThenReturn(nil, fun.ErrNotFound("not found"))
	repo := mock.Mock[ports.ActionTokenRepo]()
	ops := emailOps(t, actors, actionMgr(repo))
	err := ops.ResendVerification(context.Background(), models.ResendVerificationInput{Email: email})
	if err != nil {
		t.Fatalf("unknown email must be a no-op, got %v", err)
	}

	// already verified → 200, nothing dispatched
	mock.SetUp(t)
	actor := testActor()
	now := time.Now()
	actor.VerifiedAt = &now
	sender := mock.Mock[ports.EmailSender]()
	actors = mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByEmail(mock.AnyContext(), mock.Equal(email), mock.Any[*uuid.UUID]())).ThenReturn(&actor, nil)
	ops = newOpsWithSender(t, actors, sender)
	verifyErr := ops.ResendVerification(context.Background(), models.ResendVerificationInput{Email: email})
	if verifyErr != nil {
		t.Fatalf("verified actor must be a no-op, got %v", verifyErr)
	}
	_ = mock.Verify(sender, mock.Times(0)).SendVerify(mock.AnyContext(), mock.Any[*models.Actor](), mock.Any[*models.Project]())
}

func TestResendVerificationDispatches(t *testing.T) {
	mock.SetUp(t)
	email := "user@example.com"
	actor := testActor()

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByEmail(mock.AnyContext(), mock.Equal(email), mock.Any[*uuid.UUID]())).ThenReturn(&actor, nil)
	sender := mock.Mock[ports.EmailSender]()
	_ = mock.When(sender.SendVerify(mock.AnyContext(), mock.Equal(&actor), mock.Any[*models.Project]())).ThenReturn(nil)

	ops := newOpsWithSender(t, actors, sender)
	err := ops.ResendVerification(context.Background(), models.ResendVerificationInput{Email: email})
	if err != nil {
		t.Fatalf("ResendVerification: %v", err)
	}
	_ = mock.Verify(sender, mock.Times(1)).SendVerify(mock.AnyContext(), mock.Equal(&actor), mock.Any[*models.Project]())
}

func TestForgotPasswordSilentNoOps(t *testing.T) {
	mock.SetUp(t)
	email := "oauth@example.com"

	// OAuth-only actor (no password) → 200, nothing dispatched
	actor := testActor()
	sender := mock.Mock[ports.EmailSender]()
	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByEmail(mock.AnyContext(), mock.Equal(email), mock.Any[*uuid.UUID]())).ThenReturn(&actor, nil)
	ops := newOpsWithSender(t, actors, sender)
	err := ops.ForgotPassword(context.Background(), models.ForgotPasswordInput{Email: email})
	if err != nil {
		t.Fatalf("passwordless actor must be a no-op, got %v", err)
	}
	_ = mock.Verify(sender, mock.Times(0)).SendReset(mock.AnyContext(), mock.Any[*models.Actor](), mock.Any[*models.Project]())
}

func TestForgotPasswordDispatches(t *testing.T) {
	mock.SetUp(t)
	email := "user@example.com"
	actor := testActor()
	hash := "some-hash"
	actor.PasswordHash = &hash

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByEmail(mock.AnyContext(), mock.Equal(email), mock.Any[*uuid.UUID]())).ThenReturn(&actor, nil)
	sender := mock.Mock[ports.EmailSender]()
	_ = mock.When(sender.SendReset(mock.AnyContext(), mock.Equal(&actor), mock.Any[*models.Project]())).ThenReturn(nil)

	ops := newOpsWithSender(t, actors, sender)
	err := ops.ForgotPassword(context.Background(), models.ForgotPasswordInput{Email: email})
	if err != nil {
		t.Fatalf("ForgotPassword: %v", err)
	}
	_ = mock.Verify(sender, mock.Times(1)).SendReset(mock.AnyContext(), mock.Equal(&actor), mock.Any[*models.Project]())
}

// newOpsWithSender builds authn ops with a stubbed email sender.
func newOpsWithSender(t *testing.T, actors ports.ActorRepo, sender ports.EmailSender) *Operations {
	t.Helper()
	repo := mock.Mock[ports.ActionTokenRepo]()
	return NewOperations(
		actors,
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.PlatformRolesRepo](),
		tokens.NewManager(mock.Mock[ports.CryptoKeysRepo](), mock.Mock[ports.BlacklistRepo](), mock.Mock[ports.ActorRepo](), mock.Mock[ports.ProjectRepo](), tokens.Config{}),
		actionMgr(repo),
		sender,
	)
}
