package models

import (
	"context"
	"encoding/json"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Subject represents the full 'sub' claim object
type Subject struct {
	ID           uuid.UUID       `json:"id"`
	ProjectID    *uuid.UUID      `json:"project_id"`
	Email        *string         `json:"email"`
	Type         ActorType       `json:"type"`
	Capabilities json.RawMessage `json:"capabilities"`
	Metadata     json.RawMessage `json:"metadata"`
}

func SubjectFromAccessSub(sub AccessSub) Subject {
	return Subject{
		ID:           sub.ID,
		ProjectID:    sub.ProjectID,
		Email:        sub.Email,
		Type:         sub.Type,
		Capabilities: sub.Capabilities,
		Metadata:     sub.Metadata,
	}
}

type Credential struct {
	ID   *uuid.UUID     `json:"id"` // Applicable for stateful credentials like api keys
	Type CredentialType `json:"type"`
	Raw  string         `json:"-"`
}

type Identity struct {
	Sub  Subject    `json:"sub"`
	Cred Credential `json:"cred"`
}

type contextKey string

const identityContextKey contextKey = "identity"

func WithIdentity(ctx context.Context, identity *Identity) context.Context {
	return context.WithValue(ctx, identityContextKey, identity)
}

func RequireIdentity(ctx context.Context) (*Identity, error) {
	val := ctx.Value(identityContextKey)
	if val == nil {
		return nil, fun.ErrUnauthorized("identity not found in context")
	}
	id, ok := val.(*Identity)
	if !ok {
		return nil, fun.Errf("invalid identity type, was: %T", val).Internal()
	}
	return id, nil
}
