package authn

import (
	"context"
	"testing"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

const testHMAC = "test-hmac-secret"

func emailOps(t *testing.T, actors ports.ActorRepo, tokens ports.ActionTokenRepo) *Operations {
	t.Helper()
	return NewOperations(
		actors,
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.PlatformRolesRepo](),
		mock.Mock[ports.CryptoKeysRepo](),
		mock.Mock[ports.BlacklistRepo](),
		mock.Mock[ports.ExternalIdentitiesRepo](),
		mockOAuthProviderOps(t),
		mock.Mock[ports.OAuthLoginStatesRepo](),
		tokens,
		mock.Mock[ports.EmailSender](),
		[]byte(testHMAC),
	)
}

// mintActionToken signs an action token for the actor and purpose.
func mintActionToken(t *testing.T, actorID uuid.UUID, purpose models.ActionTokenPurpose, expiresIn time.Duration) (string, uuid.UUID) {
	t.Helper()
	jti := uuid.New()
	token, err := crypto.SignHMACJWT(mintClaims(actorID, purpose, jti, expiresIn), []byte(testHMAC))
	if err != nil {
		t.Fatalf("SignHMACJWT: %v", err)
	}
	return token, jti
}

func mintClaims(actorID uuid.UUID, purpose models.ActionTokenPurpose, jti uuid.UUID, expiresIn time.Duration) models.ActionTokenClaims {
	return models.ActionTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   actorID.String(),
			ID:        jti.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Purpose: string(purpose),
	}
}

// stubTokenRecord returns a stored action-token row for the minted jti.
func stubTokenRecord(jti uuid.UUID, purpose models.ActionTokenPurpose, actorID uuid.UUID, used bool) *models.ActionToken {
	rec := &models.ActionToken{
		JTI:       jti,
		Purpose:   purpose,
		ActorID:   actorID,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if used {
		now := time.Now()
		rec.UsedAt = &now
	}
	return rec
}

func TestVerifyEmailSuccess(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	actorID := actor.ID
	token, jti := mintActionToken(t, actorID, models.EmailVerifyActionTokenPurpose, 10*time.Minute)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.SetVerifiedAt(mock.AnyContext(), mock.Equal(actorID), mock.Any[time.Time]())).ThenReturn(nil)

	tokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(tokens.GetByJTI(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.EmailVerifyActionTokenPurpose, actorID, false), nil)
	_ = mock.When(tokens.Consume(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.EmailVerifyActionTokenPurpose, actorID, true), nil)

	ops := emailOps(t, actors, tokens)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	_, _ = mock.Verify(tokens, mock.Times(1)).Consume(mock.AnyContext(), mock.Equal(jti))
	_ = mock.Verify(actors, mock.Times(1)).SetVerifiedAt(mock.AnyContext(), mock.Equal(actorID), mock.Any[time.Time]())
}

func TestVerifyEmailIdempotentWhenAlreadyVerified(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	now := time.Now()
	actor.VerifiedAt = &now
	token, jti := mintActionToken(t, actor.ID, models.EmailVerifyActionTokenPurpose, 10*time.Minute)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)

	tokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(tokens.GetByJTI(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.EmailVerifyActionTokenPurpose, actor.ID, true), nil)

	ops := emailOps(t, actors, tokens)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err != nil {
		t.Fatalf("re-click on verified account must succeed, got %v", err)
	}
}

func TestVerifyEmailRejectsUsedTokenOnUnverifiedAccount(t *testing.T) {
	mock.SetUp(t)
	actor := testActor() // unverified
	token, jti := mintActionToken(t, actor.ID, models.EmailVerifyActionTokenPurpose, 10*time.Minute)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actor.ID))).ThenReturn(&actor, nil)

	tokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(tokens.GetByJTI(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.EmailVerifyActionTokenPurpose, actor.ID, true), nil)

	ops := emailOps(t, actors, tokens)
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for replayed token, got %v", err)
	}
}

func TestVerifyEmailRejectsWrongPurpose(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	token, _ := mintActionToken(t, actor.ID, models.PasswordResetActionTokenPurpose, 10*time.Minute)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), mock.Mock[ports.ActionTokenRepo]())
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: token})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for reset token on verify, got %v", err)
	}
}

func TestVerifyEmailRejectsGarbageToken(t *testing.T) {
	mock.SetUp(t)
	ops := emailOps(t, mock.Mock[ports.ActorRepo](), mock.Mock[ports.ActionTokenRepo]())
	err := ops.VerifyEmail(context.Background(), models.VerifyEmailInput{Token: "not-a-jwt"})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for garbage token, got %v", err)
	}
}

func TestResetPasswordSuccess(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	token, jti := mintActionToken(t, actor.ID, models.PasswordResetActionTokenPurpose, 10*time.Minute)

	actors := mock.Mock[ports.ActorRepo]()
	captor := mock.Captor[string]()
	_ = mock.When(actors.UpdatePasswordHash(mock.AnyContext(), mock.Equal(actor.ID), captor.Capture())).ThenReturn(nil)

	tokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(tokens.GetByJTI(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.PasswordResetActionTokenPurpose, actor.ID, false), nil)
	_ = mock.When(tokens.Consume(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.PasswordResetActionTokenPurpose, actor.ID, true), nil)

	ops := emailOps(t, actors, tokens)
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
	_, _ = mock.Verify(tokens, mock.Times(1)).Consume(mock.AnyContext(), mock.Equal(jti))
}

func TestResetPasswordRejectsUsedToken(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	token, jti := mintActionToken(t, actor.ID, models.PasswordResetActionTokenPurpose, 10*time.Minute)

	tokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(tokens.GetByJTI(mock.AnyContext(), mock.Equal(jti))).ThenReturn(stubTokenRecord(jti, models.PasswordResetActionTokenPurpose, actor.ID, true), nil)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), tokens)
	err := ops.ResetPassword(context.Background(), models.ResetPasswordInput{Token: token, NewPassword: "NewPassword123!"})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request for used reset token, got %v", err)
	}
}

func TestResetPasswordRejectsExpiredToken(t *testing.T) {
	mock.SetUp(t)
	actor := testActor()
	token, jti := mintActionToken(t, actor.ID, models.PasswordResetActionTokenPurpose, -time.Hour)

	tokens := mock.Mock[ports.ActionTokenRepo]()
	rec := stubTokenRecord(jti, models.PasswordResetActionTokenPurpose, actor.ID, false)
	rec.ExpiresAt = time.Now().Add(-time.Hour)
	_ = mock.When(tokens.GetByJTI(mock.AnyContext(), mock.Equal(jti))).ThenReturn(rec, nil)

	ops := emailOps(t, mock.Mock[ports.ActorRepo](), tokens)
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
	ops := emailOps(t, actors, mock.Mock[ports.ActionTokenRepo]())
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
	return NewOperations(
		actors,
		mock.Mock[ports.ProjectRepo](),
		mock.Mock[ports.PlatformRolesRepo](),
		mock.Mock[ports.CryptoKeysRepo](),
		mock.Mock[ports.BlacklistRepo](),
		mock.Mock[ports.ExternalIdentitiesRepo](),
		mockOAuthProviderOps(t),
		mock.Mock[ports.OAuthLoginStatesRepo](),
		mock.Mock[ports.ActionTokenRepo](),
		sender,
		[]byte(testHMAC),
	)
}
