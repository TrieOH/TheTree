package api_keys

import (
	"context"
	"encoding/json"
	"time"

	"IdentityX/models"
	"lib/api_keys"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

// Authenticate resolves a raw API key into the caller's identity: parse,
// prefix lookup, HMAC verification, expiry, actor resolution, and
// capability shaping all happen here, behind the same seam the minting
// side (Create) crosses. The auth middleware calls this instead of
// reaching into the repos, so key verification and key minting live in
// one module and the capability format (resource:action pairs on the
// identity's subject) has a single home.
func (o *Operations) Authenticate(ctx context.Context, rawKey string) (*models.Identity, error) {
	ctx, span := telemetry.StartSpan(ctx, "Authenticate")
	defer span.End()

	key, err := api_keys.ParseAPIKey(rawKey)
	if err != nil {
		telemetry.Log().Warn("api key parse failed", zap.Error(err))
		return nil, fun.ErrForbidden("invalid api key")
	}

	apiKey, err := o.apiKeys.GetByPrefix(ctx, key.DisplayPrefix)
	if err != nil {
		telemetry.Log().Warn("api key lookup failed", zap.String("prefix", key.DisplayPrefix), zap.Error(err))
		return nil, fun.ErrForbidden("invalid api key")
	}

	if !api_keys.VerifyAPIKey(rawKey, apiKey.KeyHash, o.hmacSecret) {
		telemetry.Log().Warn("api key hmac mismatch", zap.String("prefix", key.DisplayPrefix))
		return nil, fun.ErrForbidden("invalid api key")
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		telemetry.Log().Warn("api key expired", zap.String("prefix", key.DisplayPrefix), zap.Time("expires_at", *apiKey.ExpiresAt))
		return nil, fun.ErrForbidden("api key expired")
	}

	actor, err := o.actors.GetByID(ctx, apiKey.SubjectID)
	if err != nil {
		telemetry.Log().Warn("actor lookup failed", zap.String("prefix", key.DisplayPrefix), zap.String("subject_id", apiKey.SubjectID.String()), zap.Error(err))
		return nil, fun.ErrForbidden("invalid api key")
	}

	caps, err := o.capabilities.ListByAPIKeyPrefix(ctx, key.DisplayPrefix)
	if err != nil {
		return nil, err
	}
	pairs := make([]string, len(caps))
	for i, c := range caps {
		pairs[i] = c.Resource + ":" + c.Action
	}
	capJSON, _ := json.Marshal(pairs)

	telemetry.Log().Debug("api key authenticated", zap.String("prefix", key.DisplayPrefix), zap.Strings("capabilities", pairs))
	return &models.Identity{
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
	}, nil
}
