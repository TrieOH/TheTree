package app

import (
	"context"
	"encoding/json"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/api_keys"
	"lib/crypto"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	mws "github.com/MintzyG/fun/middlewares"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (app *IdentityX) SetupAuthMiddlewares(
	cryptoKeysRepo ports.CryptoKeysRepo,
	apiKeysRepo ports.APIKeysRepo,
	actorsRepo ports.ActorRepo,
	capabilitiesRepo ports.CapabilityRepo,
	blacklistRepo ports.BlacklistRepo,
) *mws.Middleware[*models.AccessClaims] {
	keyFunc := func(ctx context.Context, tokenStr string) (*models.AccessClaims, error) {
		claims := &models.AccessClaims{}
		token, err := crypto.OpenUnverified(tokenStr, claims)
		if err != nil {
			return nil, err
		}
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fun.ErrUnauthorized("missing kid")
		}
		keyID, err := uuid.Parse(kid)
		if err != nil {
			return nil, fun.ErrUnauthorized("invalid kid")
		}
		cryptoKey, err := cryptoKeysRepo.GetByID(ctx, keyID)
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrUnauthorized("outdated token")
		}
		if err != nil {
			return nil, err
		}
		if cryptoKey.Status == "revoked" {
			return nil, fun.ErrUnauthorized("token signing key revoked")
		}
		_, err = crypto.VerifyToken(tokenStr, cryptoKey.PublicKey, claims)
		if err != nil {
			return nil, fun.ErrUnauthorized("invalid access token")
		}
		// a token blacklisted at logout (or by refresh rotation) must not
		// authenticate anymore
		_, err = blacklistRepo.GetByTargetAndType(ctx, claims.ID, models.BlacklistEntryTypeToken)
		if err == nil {
			return nil, fun.ErrUnauthorized("token has been revoked")
		}
		if !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		return claims, nil
	}

	jwtHook := func(ctx context.Context, claims *models.AccessClaims) (context.Context, error) {
		identity := &models.Identity{
			Sub: models.SubjectFromAccessSub(&claims.Sub),
			Cred: models.Credential{
				Type: models.TokenCredentialType,
			},
		}
		return models.WithIdentity(ctx, identity), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		key, err := api_keys.ParseAPIKey(rawKey)
		if err != nil {
			telemetry.Log().Warn("api key parse failed", zap.Error(err))
			return nil, fun.ErrForbidden("invalid api key")
		}

		apiKey, err := apiKeysRepo.GetByPrefix(ctx, key.DisplayPrefix)
		if err != nil {
			telemetry.Log().Warn("api key lookup failed", zap.String("prefix", key.DisplayPrefix), zap.Error(err))
			return nil, fun.ErrForbidden("invalid api key")
		}

		if !api_keys.VerifyAPIKey(rawKey, apiKey.KeyHash, []byte(app.cfg.HmacSecret)) {
			telemetry.Log().Warn("api key hmac mismatch", zap.String("prefix", key.DisplayPrefix))
			return nil, fun.ErrForbidden("invalid api key")
		}

		if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
			telemetry.Log().Warn("api key expired", zap.String("prefix", key.DisplayPrefix), zap.Time("expires_at", *apiKey.ExpiresAt))
			return nil, fun.ErrForbidden("api key expired")
		}

		actor, err := actorsRepo.GetByID(ctx, apiKey.SubjectID)
		if err != nil {
			telemetry.Log().Warn("actor lookup failed", zap.String("prefix", key.DisplayPrefix), zap.String("subject_id", apiKey.SubjectID.String()), zap.Error(err))
			return nil, fun.ErrForbidden("invalid api key")
		}

		caps, err := capabilitiesRepo.ListByAPIKeyPrefix(ctx, key.DisplayPrefix)
		if err != nil {
			return nil, err
		}
		pairs := make([]string, len(caps))
		for i, c := range caps {
			pairs[i] = c.Resource + ":" + c.Action
		}
		capJSON, _ := json.Marshal(pairs)

		ctx = models.WithIdentity(ctx, &models.Identity{
			Sub: models.Subject{
				ID:           actor.ID,
				ProjectID:    actor.ProjectID,
				Email:        actor.Email,
				Type:         actor.Type,
				Capabilities: capJSON,
				Metadata:     actor.Metadata,
			},
			Cred: models.Credential{
				ID:   &apiKey.ID,
				Type: models.APIKeyCredentialType,
				Raw:  rawKey,
			},
		})
		telemetry.Log().Debug("api key authenticated", zap.String("prefix", key.DisplayPrefix), zap.Strings("capabilities", pairs))
		return ctx, nil
	}
	return mws.New[*models.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
