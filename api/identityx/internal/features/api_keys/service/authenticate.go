package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"IdentityX/errs"
	"IdentityX/models"
	"lib/api_keys"
	"lib/telemetry"
)

// Authenticate resolves a raw API key into the caller's identity: parse,
// prefix lookup, HMAC verification, expiry, actor resolution, and
// capability shaping all happen here, behind the same seam the minting
// side (Create) crosses. The auth middleware calls this instead of
// reaching into the repos, so key verification and key minting live in
// one module and the capability format (resource:action pairs on the
// identity's subject) has a single home.
func (o *Operations) Authenticate(ctx context.Context, rawKey string) (*models.Identity, error) {
	key, err := api_keys.ParseAPIKey(rawKey)
	if err != nil {
		telemetry.Log().Warn("api key parse failed", telemetry.Err(err))
		return nil, errs.ErrInvalidAPIKey
	}

	apiKey, err := o.APIKeys.GetByPrefix(ctx, key.DisplayPrefix)
	if err != nil {
		telemetry.Log().Warn("api key lookup failed", slog.String("prefix", key.DisplayPrefix), telemetry.Err(err))
		return nil, errs.ErrInvalidAPIKey
	}

	if !api_keys.VerifyAPIKey(rawKey, apiKey.KeyHash, o.HmacSecret) {
		telemetry.Log().Warn("api key hmac mismatch", slog.String("prefix", key.DisplayPrefix))
		return nil, errs.ErrInvalidAPIKey
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		telemetry.Log().Warn("api key expired", slog.String("prefix", key.DisplayPrefix), slog.Time("expires_at", *apiKey.ExpiresAt))
		return nil, errs.ErrAPIKeyExpired
	}

	actor, err := o.Actors.GetByID(ctx, apiKey.SubjectID)
	if err != nil {
		telemetry.Log().Warn("actor lookup failed", slog.String("prefix", key.DisplayPrefix), slog.String("subject_id", apiKey.SubjectID.String()), telemetry.Err(err))
		return nil, errs.ErrInvalidAPIKey
	}

	caps, err := o.Capabilities.ListByAPIKeyPrefix(ctx, key.DisplayPrefix)
	if err != nil {
		return nil, err
	}
	pairs := make([]string, len(caps))
	for i, c := range caps {
		pairs[i] = c.Resource + ":" + c.Action
	}
	capJSON, _ := json.Marshal(pairs)

	telemetry.Log().Debug("api key authenticated", slog.String("prefix", key.DisplayPrefix), slog.Any("capabilities", pairs))
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
