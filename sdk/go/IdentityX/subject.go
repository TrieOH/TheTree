package idx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CredentialType string

const (
	TokenCredentialType  CredentialType = "token"
	APIKeyCredentialType CredentialType = "api_key"
)

// Subject represents the full 'sub' claim object
type Subject struct {
	ID           uuid.UUID       `json:"id"`
	ProjectID    *uuid.UUID      `json:"project_id"`
	Email        *string         `json:"email"`
	Type         ActorType       `json:"type"`
	VerifiedAt   *time.Time      `json:"verified_at"`
	Capabilities json.RawMessage `json:"capabilities"`
	Metadata     json.RawMessage `json:"metadata"`
}

func SubjectFromAccessSub(sub AccessSub) Subject {
	return Subject(sub)
}

type Credential struct {
	ID   *uuid.UUID     `json:"id"` // Applicable for stateful credentials like api keys
	Type CredentialType `json:"type"`
	Raw  string         `json:"-"`
}

type Identity struct {
	Sub  Subject    `json:"subject"`
	Cred Credential `json:"cred"`
}

type contextKey string

const identityContextKey contextKey = "identity"

func WithIdentity(ctx context.Context, identity *Identity) context.Context {
	return context.WithValue(ctx, identityContextKey, identity)
}

// ErrNoIdentity is RequireIdentity on a request no auth middleware
// authenticated. The Backends' error maps answer it with a 401.
var ErrNoIdentity = errors.New("identity not found in context")

// RequireIdentity returns the caller's identity, or ErrNoIdentity.
func RequireIdentity(ctx context.Context) (*Identity, error) {
	val := ctx.Value(identityContextKey)
	if val == nil {
		return nil, ErrNoIdentity
	}
	id, ok := val.(*Identity)
	if !ok {
		return nil, fmt.Errorf("idx: invalid identity type %T in context", val)
	}
	return id, nil
}
