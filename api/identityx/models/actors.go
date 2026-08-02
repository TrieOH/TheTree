package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuthMethod string

const (
	PasswordAuthMethod AuthMethod = "password"
	APIKeyAuthMethod   AuthMethod = "api_key"
	GoogleAuthMethod   AuthMethod = "google_auth"
	GithubAuthMethod   AuthMethod = "github_auth"
)

type ActorType string

const (
	HumanActorType   ActorType = "human"
	ServiceActorType ActorType = "service"
	MachineActorType ActorType = "machine"
)

type Actor struct {
	ID           uuid.UUID        `json:"id"`
	ProjectID    *uuid.UUID       `json:"project_id"`
	AuthMethod   AuthMethod       `json:"auth_method"`
	VerifiedAt   *time.Time       `json:"verified_at"`
	PasswordHash *string          `json:"-"`
	Email        *string          `json:"email"`
	Type         ActorType        `json:"type"`
	Metadata     *json.RawMessage `json:"metadata"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    *time.Time       `json:"deleted_at"`
}

// ActorProfile is a single profile document per actor. SchemaVersion points
// at the profile schema version the document was validated against;
// Outdated flags documents that no longer validate against the active
// schema and were not migrated (admin resolves them manually).
type ActorProfile struct {
	ActorID       uuid.UUID       `json:"actor_id"`
	Profile       json.RawMessage `json:"profile"`
	SchemaVersion int             `json:"schema_version"`
	Outdated      bool            `json:"outdated"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type OAuthProvider string

const (
	GoogleIdentityProvider OAuthProvider = "google"
	GithubIdentityProvider OAuthProvider = "github"
)

type ActorExternalIdentities struct {
	ID                    uuid.UUID     `json:"id"`
	ActorID               uuid.UUID     `json:"actor_id"`
	Provider              OAuthProvider `json:"provider"`
	Subject               string        `json:"subject"`
	Email                 *string       `json:"email"`
	EncryptedAccessToken  *string       `json:"-"`
	EncryptedRefreshToken *string       `json:"-"`
	TokenExpiresAt        *time.Time    `json:"token_expires_at"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

type CreateActorInput struct {
	ProjectID  *uuid.UUID `json:"project_id"`
	AuthMethod AuthMethod `json:"auth_method" validate:"required,oneof=password api_key"`
	Type       ActorType  `json:"type"        validate:"required,oneof=human service machine"`
	Email      *string    `json:"email"`
}
