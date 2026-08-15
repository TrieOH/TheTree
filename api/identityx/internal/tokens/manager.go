// Package tokens owns the token lifecycle for IdentityX: opening,
// verifying, minting, rotating, and revoking signed tokens at one
// interface. The auth middleware (verify), login and OAuth callback
// (mint), refresh (rotate), and logout (revoke) cross the same seam, so
// key resolution, blacklist shaping, and entry reasons live in one module
// instead of being re-derived per caller.
package tokens

import (
	"context"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Config carries the token-policy knobs the Manager needs. Resolved once
// at construction (wire time), not read from the environment per call.
type Config struct {
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// Option configures a Manager at construction.
type Option func(*Manager)

// WithClock overrides the Manager's clock, so tests can pin exact token
// lifetimes instead of asserting "non-empty and different".
func WithClock(now func() time.Time) Option {
	return func(m *Manager) { m.now = now }
}

// Manager owns the token lifecycle. Rotation and revocation are
// fail-closed: if a blacklist append fails, the operation errors instead
// of reporting success while the old pair stays alive.
type Manager struct {
	cryptoKeys ports.CryptoKeysRepo
	blacklist  ports.BlacklistRepo
	actors     ports.ActorRepo
	projects   ports.ProjectRepo
	cfg        Config
	now        func() time.Time
}

func NewManager(cryptoKeys ports.CryptoKeysRepo, blacklist ports.BlacklistRepo, actors ports.ActorRepo, projects ports.ProjectRepo, cfg Config, opts ...Option) *Manager {
	m := &Manager{
		cryptoKeys: cryptoKeys,
		blacklist:  blacklist,
		actors:     actors,
		projects:   projects,
		cfg:        cfg,
		now:        time.Now,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Verify opens tokenStr, resolves the signing key by kid (revoked keys are
// rejected), verifies the signature, and rejects blacklisted tokens. The
// parsed claims are populated into claims.
func (m *Manager) Verify(ctx context.Context, tokenStr string, claims jwt.Claims) error {
	token, err := crypto.OpenUnverified(tokenStr, claims)
	if err != nil {
		return err
	}

	key, err := m.keyForToken(ctx, token)
	if err != nil {
		return err
	}

	_, err = crypto.VerifyToken(tokenStr, key.PublicKey, claims)
	if err != nil {
		return fun.ErrUnauthorized("invalid token")
	}

	return m.ensureNotBlacklisted(ctx, tokenJTI(claims))
}

// Mint signs a fresh access/refresh pair for the actor with the active
// signing key of the actor's scope. The refresh token carries the access
// token's jti so a later rotation can revoke the access token along with
// the refresh token that anchors it.
func (m *Manager) Mint(ctx context.Context, actor *models.Actor) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "tokens.Mint")
	defer span.End()

	activeKeyPair, err := m.cryptoKeys.GetActive(ctx, models.SigningCryptoKeyType, actor.ProjectID)
	if err != nil {
		return nil, err
	}

	now := m.now()
	accessJTI := uuid.New()
	refreshJTI := uuid.New()
	accessExpiresAt := now.Add(m.cfg.AccessTTL)
	refreshExpiresAt := now.Add(m.cfg.RefreshTTL)

	accessPayload, err := m.newAccessToken(*actor, accessJTI, activeKeyPair.ID, accessExpiresAt, now)
	if err != nil {
		return nil, err
	}
	refreshPayload, err := m.newRefreshToken(actor, refreshJTI, accessJTI, activeKeyPair.ID, refreshExpiresAt, now)
	if err != nil {
		return nil, err
	}

	kp := activeKeyPair.ToKeyPair()
	accessToken, err := crypto.SignToken(accessPayload, kp)
	if err != nil {
		return nil, err
	}
	refreshToken, err := crypto.SignToken(refreshPayload, kp)
	if err != nil {
		return nil, err
	}

	return &models.UserTokensOutput{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		Domain:           m.cfg.Issuer,
	}, nil
}

// Rotate redeems a refresh token for a fresh pair: it verifies the old
// refresh token, blacklists it and the access token it anchors (the jti
// linkage), re-reads the actor, and mints a new pair. Fail-closed: if
// either blacklist append fails, the old pair survives and the rotation
// errors instead of minting on top of it.
func (m *Manager) Rotate(ctx context.Context, refreshToken string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "tokens.Rotate")
	defer span.End()

	refreshClaims := &models.RefreshClaims{}
	err := m.Verify(ctx, refreshToken, refreshClaims)
	if err != nil {
		return nil, err
	}

	err = m.blacklist.Append(ctx, m.tokenEntry(refreshClaims, refreshClaims.ID, "refresh"))
	if err != nil {
		return nil, err
	}
	err = m.blacklist.Append(ctx, m.tokenEntry(refreshClaims, refreshClaims.Sub.AccessJTI.String(), "refresh"))
	if err != nil {
		return nil, err
	}

	actor, err := m.actors.GetByID(ctx, refreshClaims.Sub.ID)
	if err != nil {
		return nil, err
	}
	return m.Mint(ctx, actor)
}

// Revoke ends a session: it blacklists the access token and, when the
// refresh token is still verifiable, the refresh token too. The access
// token is deliberately opened without verifying — a dead token must not
// fail the logout, but its signing key is needed to check the refresh
// token. Identity for the blacklist entries comes from the claims (the
// middleware built its identity from the same claims). Fail-closed: a
// blacklist append that fails surfaces as an error rather than a silent
// "logged out".
func (m *Manager) Revoke(ctx context.Context, accessToken, refreshToken string) error {
	ctx, span := telemetry.StartSpan(ctx, "tokens.Revoke")
	defer span.End()

	accessClaims := &models.AccessClaims{}
	token, err := crypto.OpenUnverified(accessToken, accessClaims)
	if err != nil {
		return fun.ErrUnauthorized("invalid access token")
	}
	key, err := m.keyForToken(ctx, token)
	if err != nil {
		return err
	}

	err = m.blacklist.Append(ctx, m.tokenEntry(accessClaims, accessClaims.ID, "logout"))
	if err != nil {
		return err
	}

	// A dead refresh token (expired, revoked, garbage) must not fail the
	// logout: the access token is already blacklisted, so the session is
	// over either way. Only a verified refresh token gets blacklisted.
	refreshClaims := &models.RefreshClaims{}
	_, err = crypto.VerifyToken(refreshToken, key.PublicKey, refreshClaims)
	if err != nil {
		return nil
	}
	return m.blacklist.Append(ctx, m.tokenEntry(refreshClaims, refreshClaims.ID, "logout"))
}

// JWKS publishes the public signing keys of a scope as a JWKS document, so
// clients can verify tokens without holding keys themselves. A missing
// project surfaces as not found.
func (m *Manager) JWKS(ctx context.Context, projectID *uuid.UUID) (map[string]any, error) {
	ctx, span := telemetry.StartSpan(ctx, "tokens.JWKS")
	defer span.End()

	if projectID != nil {
		_, err := m.projects.GetByID(ctx, *projectID)
		if err != nil {
			return nil, err
		}
	}

	keys, err := m.cryptoKeys.GetActiveSigningKeys(ctx, projectID)
	if err != nil {
		return nil, err
	}

	jwks := make([]map[string]any, 0, len(keys))
	for _, k := range keys {
		jwk, err := crypto.PublicKeyToJWKS(k.ID.String(), k.PublicKey)
		if err != nil {
			telemetry.Log().Warn("skipping malformed key", zap.String("key_id", k.ID.String()), zap.Error(err))
			continue
		}
		jwks = append(jwks, jwk)
	}

	return map[string]any{"keys": jwks}, nil
}

// tokenEntry shapes the blacklist entry that revokes a token. The reason
// distinguishes rotation ("refresh") from logout ("logout").
func (m *Manager) tokenEntry(claims any, target, reason string) models.BlacklistEntry {
	switch c := claims.(type) {
	case *models.RefreshClaims:
		return models.BlacklistEntry{
			CreatedByActorID: &c.Sub.ID,
			ProjectID:        c.Sub.ProjectID,
			Type:             models.BlacklistEntryTypeToken,
			Target:           target,
			Reason:           &reason,
			ExpiresAt:        &c.ExpiresAt.Time,
		}
	case *models.AccessClaims:
		return models.BlacklistEntry{
			CreatedByActorID: &c.Sub.ID,
			ProjectID:        c.Sub.ProjectID,
			Type:             models.BlacklistEntryTypeToken,
			Target:           target,
			Reason:           &reason,
			ExpiresAt:        &c.ExpiresAt.Time,
		}
	default:
		return models.BlacklistEntry{}
	}
}

// keyForToken resolves the signing key for an already-opened token,
// rejecting missing and revoked keys.
func (m *Manager) keyForToken(ctx context.Context, token *jwt.Token) (*models.CryptoKey, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fun.ErrUnauthorized("missing kid")
	}
	keyID, err := uuid.Parse(kid)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid kid")
	}
	key, err := m.cryptoKeys.GetByID(ctx, keyID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("outdated token")
	}
	if err != nil {
		return nil, err
	}
	if key.Status == "revoked" {
		return nil, fun.ErrUnauthorized("token signing key revoked")
	}
	return key, nil
}

// ensureNotBlacklisted rejects tokens revoked at logout or refresh
// rotation. A missing entry (never revoked) passes.
func (m *Manager) ensureNotBlacklisted(ctx context.Context, jti string) error {
	_, err := m.blacklist.GetByTargetAndType(ctx, jti, models.BlacklistEntryTypeToken)
	if err == nil {
		return fun.ErrUnauthorized("token has been revoked")
	}
	if !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	return nil
}

// tokenJTI extracts the jti (RegisteredClaims.ID) from the two claim types
// IdentityX signs; the jwt.Claims interface does not expose it.
func tokenJTI(claims jwt.Claims) string {
	switch c := claims.(type) {
	case *models.AccessClaims:
		return c.ID
	case *models.RefreshClaims:
		return c.ID
	default:
		return ""
	}
}

// newAccessToken builds the signing string (header.payload) of an access
// token for the actor, stamped with the given jti, kid, and lifetime.
func (m *Manager) newAccessToken(actor models.Actor, jti, kid uuid.UUID, expiresAt, now time.Time) ([]byte, error) {
	claims := models.AccessClaims{
		Sub: models.AccessSub{
			ID:         actor.ID,
			ProjectID:  actor.ProjectID,
			Email:      actor.Email,
			Type:       actor.Type,
			VerifiedAt: actor.VerifiedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    m.cfg.Issuer,
			ID:        jti.String(),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

// newRefreshToken builds the signing string of a refresh token anchored to
// the access token's jti, so rotation can revoke both with one claims set.
func (m *Manager) newRefreshToken(actor *models.Actor, jti, accessJTI, kid uuid.UUID, expiresAt, now time.Time) ([]byte, error) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSub{
			ID:        actor.ID,
			ProjectID: actor.ProjectID,
			AccessJTI: accessJTI,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    m.cfg.Issuer,
			ID:        jti.String(),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}
