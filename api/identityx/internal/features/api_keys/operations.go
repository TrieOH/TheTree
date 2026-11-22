// Package api_keys is the api keys feature: its repo, handlers and
// operations. Build wires it.
package api_keys

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"encoding/json"
	"fmt"
	"lib/apikeys"
	"lib/telemetry"
	"log/slog"
	"time"
	"uuid"

	libauthz "lib/authz"
)

type Deps struct {
	HmacSecret   []byte
	Actors       ports.ActorRepo
	APIKeys      ports.APIKeysRepo
	Capabilities ports.CapabilityRepo
	Projects     ports.ProjectRepo
	Authz        *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// Authenticate resolves a raw API key into the caller's identity: parse,
// prefix lookup, HMAC verification, expiry, actor resolution, and
// capability shaping all happen here, behind the same seam the minting
// side (Create) crosses. The auth middleware calls this instead of
// reaching into the repos, so key verification and key minting live in
// one module and the capability format (resource:action pairs on the
// identity's subject) has a single home.
func (o *Operations) Authenticate(ctx context.Context, rawKey string) (*models.Identity, error) {
	key, err := apikeys.ParseAPIKey(rawKey)
	if err != nil {
		telemetry.Log().WarnContext(ctx, "api key parse failed", telemetry.Err(err))
		return nil, errs.ErrInvalidAPIKey
	}

	apiKey, err := o.APIKeys.GetByPrefix(ctx, key.DisplayPrefix)
	if err != nil {
		telemetry.Log().WarnContext(ctx, "api key lookup failed", slog.String("prefix", key.DisplayPrefix), telemetry.Err(err))
		return nil, errs.ErrInvalidAPIKey
	}

	if !apikeys.VerifyAPIKey(rawKey, apiKey.KeyHash, o.HmacSecret) {
		telemetry.Log().WarnContext(ctx, "api key hmac mismatch", slog.String("prefix", key.DisplayPrefix))
		return nil, errs.ErrInvalidAPIKey
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		telemetry.Log().WarnContext(ctx, "api key expired", slog.String("prefix", key.DisplayPrefix), slog.Time("expires_at", *apiKey.ExpiresAt))
		return nil, errs.ErrAPIKeyExpired
	}

	actor, err := o.Actors.GetByID(ctx, apiKey.SubjectID)
	if err != nil {
		telemetry.Log().WarnContext(ctx, "actor lookup failed", slog.String("prefix", key.DisplayPrefix), slog.String("subject_id", apiKey.SubjectID.String()), telemetry.Err(err))
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

	telemetry.Log().DebugContext(ctx, "api key authenticated", slog.String("prefix", key.DisplayPrefix), slog.Any("capabilities", pairs))
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

func (o *Operations) Create(ctx context.Context, payload models.CreateAPIKeyInput) (*models.APIKey, string, error) {
	project, err := o.Projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		telemetry.Log().InfoContext(ctx, "Create api key", slog.String("project_id", payload.ProjectID.String()))
		return nil, "", err
	}

	created, generated, err := o.createInternal(ctx, *project, payload)
	if err != nil {
		return nil, "", err
	}

	return created, generated.Raw, nil
}

// requireKeySubject limits whom a key can be minted for besides the caller:
// a service or machine actor of the same project. A key acts as its subject,
// so minting one for anyone else (a platform admin, a human, an actor of
// another project) would let the caller act as them.
func (o *Operations) requireKeySubject(ctx context.Context, subjectID, projectID uuid.UUID) error {
	subject, err := o.Actors.GetByID(ctx, subjectID)
	if err != nil {
		return err
	}
	if subject.ProjectID == nil || *subject.ProjectID != projectID {
		return errs.ErrActorOutsideProject
	}
	if subject.Type == models.HumanActorType {
		return libauthz.ErrInsufficientRole
	}
	return nil
}

func (o *Operations) createInternal(ctx context.Context, project models.Project, payload models.CreateAPIKeyInput) (*models.APIKey, *apikeys.GeneratedAPIKey, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Every mint crosses the Access-check seam, self-bound or not: minting
	// a key for a project demands the caller hold at least admin on it.
	// CheckProject's org fallback already admits org owners with no
	// project_members row, so no legitimate path is lost.
	err = o.Authz.CheckProject(ctx, project.ID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, nil, err
	}

	actorID := ident.Sub.ID
	if payload.SubjectID != nil && *payload.SubjectID != ident.Sub.ID {
		err = o.requireKeySubject(ctx, *payload.SubjectID, project.ID)
		if err != nil {
			return nil, nil, err
		}
		actorID = *payload.SubjectID
	}

	if len(payload.Capabilities) > 0 {
		valid, err := o.Capabilities.ValidateCapabilities(ctx, &project.ID, payload.Capabilities)
		if err != nil {
			return nil, nil, err
		}
		if !valid {
			return nil, nil, errs.ErrInvalidCapabilities
		}
	}

	var generated *apikeys.GeneratedAPIKey
	generated, err = apikeys.GenerateAPIKey(project.BrandSlug, payload.Env, o.HmacSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("generate api key: %w", err)
	}

	apiKey := models.APIKey{
		SubjectID:     actorID,
		Name:          payload.Name,
		DisplayPrefix: generated.DisplayPrefix,
		KeyHash:       generated.Hash,
		ExpiresAt:     payload.ExpiresAt,
		CreatedBy:     ident.Sub.ID,
	}

	var created *models.APIKey
	created, err = o.APIKeys.Create(ctx, apiKey)
	if err != nil {
		return nil, nil, err
	}

	if len(payload.Capabilities) > 0 {
		err = o.Capabilities.AssignToAPIKey(ctx, created.ID, payload.Capabilities, ident.Sub.ID)
		if err != nil {
			return nil, nil, err
		}
	}

	return created, generated, nil
}
