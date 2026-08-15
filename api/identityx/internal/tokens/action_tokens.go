package tokens

import (
	"context"
	"errors"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ActionTokenConfig carries the per-purpose lifetimes of single-use action
// tokens. Resolved once at construction (wire time), exactly like Config.
type ActionTokenConfig struct {
	VerifyTTL time.Duration // email verification links
	ResetTTL  time.Duration // password reset links
}

// Redeem failure classes, attached to the returned AppError so callers can
// branch with errors.Is while the HTTP layer still maps them as BadRequest.
var (
	// ErrActionTokenInvalid covers unparseable, wrongly-purposed, or
	// unknown tokens. The failure always surfaces as the same generic
	// message so the endpoint never leaks token internals.
	ErrActionTokenInvalid = errors.New("action token invalid")
	// ErrActionTokenExpired covers tokens whose anti-replay record has
	// expired.
	ErrActionTokenExpired = errors.New("action token expired")
	// ErrActionTokenUsed covers tokens already consumed — Redeem's strict
	// answer; callers apply their own actor-state policy on top (the
	// verify email's idempotent fall-through is authn's business, not
	// this module's).
	ErrActionTokenUsed = errors.New("action token already used")
)

// ActionTokenOption configures an ActionTokenManager at construction.
type ActionTokenOption func(*ActionTokenManager)

// WithActionTokenClock overrides the ActionTokenManager's clock, so tests
// can pin exact token lifetimes (mirrors WithClock on Manager).
func WithActionTokenClock(now func() time.Time) ActionTokenOption {
	return func(m *ActionTokenManager) { m.now = now }
}

// ActionTokenManager owns the single-use token lifecycle: minting (sign +
// persist the anti-replay record), redeeming (purpose scope + consume),
// and sweeping expired records. emails.Sender crosses it to mint the link
// tokens, authn crosses it to redeem them, and the cleanup job crosses it
// to keep the anti-replay table bounded. The HMAC secret, the anti-replay
// repo, and the purpose→TTL policy live here — not in the callers.
type ActionTokenManager struct {
	actionTokens ports.ActionTokenRepo
	hmacSecret   []byte
	cfg          ActionTokenConfig
	now          func() time.Time
}

func NewActionTokenManager(actionTokens ports.ActionTokenRepo, hmacSecret []byte, cfg ActionTokenConfig, opts ...ActionTokenOption) *ActionTokenManager {
	m := &ActionTokenManager{
		actionTokens: actionTokens,
		hmacSecret:   hmacSecret,
		cfg:          cfg,
		now:          time.Now,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Mint signs a single-use token for the purpose, persists its anti-replay
// record, and returns the token plus the lifetime used (the caller renders
// the lifetime into the email copy). projectID distinguishes project actors
// from platform actors; it also travels on the link URL, but the claim is
// the source of truth at redemption.
func (m *ActionTokenManager) Mint(ctx context.Context, purpose models.ActionTokenPurpose, actorID uuid.UUID, projectID *uuid.UUID) (string, time.Duration, error) {
	ctx, span := telemetry.StartSpan(ctx, "tokens.ActionTokenMint")
	defer span.End()

	ttl, err := m.ttlFor(purpose)
	if err != nil {
		return "", 0, err
	}

	jti := uuid.New()
	now := m.now()
	expiresAt := now.Add(ttl)
	token, err := crypto.SignHMACJWT(models.ActionTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   actorID.String(),
			ID:        jti.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Purpose:   string(purpose),
		ProjectID: projectID,
	}, m.hmacSecret)
	if err != nil {
		return "", 0, err
	}

	_, err = m.actionTokens.Insert(ctx, models.ActionToken{
		JTI:       jti,
		Purpose:   purpose,
		ActorID:   actorID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", 0, err
	}

	return token, ttl, nil
}

// Redeem validates the HMAC signature, scopes the purpose, and atomically
// consumes the token's anti-replay record. Strict: a consumed or expired
// token is an error, with its failure class attached
// (ErrActionTokenUsed / ErrActionTokenExpired / ErrActionTokenInvalid) so
// callers can branch with errors.Is. The parsed actorID is returned even
// on used/expired errors — the signature verified, so the caller already
// knows who the token was for and can apply actor-state policy.
func (m *ActionTokenManager) Redeem(ctx context.Context, purpose models.ActionTokenPurpose, tokenStr string) (uuid.UUID, error) {
	ctx, span := telemetry.StartSpan(ctx, "tokens.ActionTokenRedeem")
	defer span.End()

	claims := &models.ActionTokenClaims{}
	_, err := crypto.ParseHMACJWT(tokenStr, claims, m.hmacSecret)
	if err != nil {
		return uuid.Nil, m.redeemError(ErrActionTokenInvalid, "invalid or expired token")
	}
	if claims.Purpose != string(purpose) {
		return uuid.Nil, m.redeemError(ErrActionTokenInvalid, "invalid or expired token")
	}
	actorID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, m.redeemError(ErrActionTokenInvalid, "invalid or expired token")
	}
	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return uuid.Nil, m.redeemError(ErrActionTokenInvalid, "invalid or expired token")
	}

	record, err := m.actionTokens.GetByJTI(ctx, jti)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return uuid.Nil, m.redeemError(ErrActionTokenInvalid, "invalid or expired token")
		}
		return uuid.Nil, err
	}
	if record.UsedAt != nil {
		return actorID, m.redeemError(ErrActionTokenUsed, "token already used")
	}
	if m.now().After(record.ExpiresAt) {
		return actorID, m.redeemError(ErrActionTokenExpired, "token expired")
	}
	_, err = m.actionTokens.Consume(ctx, jti)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			// consumed concurrently; the token is dead either way.
			return actorID, m.redeemError(ErrActionTokenUsed, "token already used")
		}
		return uuid.Nil, err
	}
	return actorID, nil
}

// DeleteExpired sweeps consumed/expired records so the anti-replay table
// stays bounded; the periodic cleanup job crosses it.
func (m *ActionTokenManager) DeleteExpired(ctx context.Context) error {
	ctx, span := telemetry.StartSpan(ctx, "tokens.ActionTokenDeleteExpired")
	defer span.End()
	return m.actionTokens.DeleteExpired(ctx)
}

// redeemError wraps a failure class with the wire-facing BadRequest error.
func (m *ActionTokenManager) redeemError(class error, msg string) error {
	return fun.Err(msg).WithErr(class).BadRequest()
}

// ttlFor resolves the lifetime for a purpose from the construction-time
// config. An unknown purpose is a programming error, not a runtime one.
func (m *ActionTokenManager) ttlFor(purpose models.ActionTokenPurpose) (time.Duration, error) {
	switch purpose {
	case models.EmailVerifyActionTokenPurpose:
		return m.cfg.VerifyTTL, nil
	case models.PasswordResetActionTokenPurpose:
		return m.cfg.ResetTTL, nil
	default:
		return 0, fun.ErrInternal("unknown action token purpose: " + string(purpose))
	}
}
