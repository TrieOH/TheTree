package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID               uuid.UUID       `json:"id"`
	OrganizationID   *uuid.UUID      `json:"organization_id"`
	OwnerID          uuid.UUID       `json:"owner_id"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Domain           *string         `json:"domain"`
	DomainVerifiedAt *time.Time      `json:"domain_verified_at"`
	Metadata         json.RawMessage `json:"metadata"`
	CreatedAt        *time.Time      `json:"created_at"`
	DeletedAt        *time.Time      `json:"deleted_at"`
}

type ProjectDomainChallenges struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Domain     string    `json:"domain"`
	Token      string    `json:"token"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	VerifiedAt time.Time `json:"verified_at"`
}

type ProjectRole string

const (
	ProjectRoleMember ProjectRole = "member"
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleOwner  ProjectRole = "owner"
)

type ProjectMembers struct {
	ProjectID uuid.UUID       `json:"project_id"`
	ActorID   uuid.UUID       `json:"actor_id"`
	Role      ProjectRole     `json:"role"`
	Metadata  json.RawMessage `json:"metadata"`
	JoinedAt  *time.Time      `json:"joined_at"`
}

type ProjectOAuthProviders struct {
	ID                    uuid.UUID     `json:"id"`
	ProjectID             uuid.UUID     `json:"project_id"`
	Provider              OAuthProvider `json:"provider"`
	ClientID              string        `json:"client_id"`
	EncryptedClientSecret string        `json:"encrypted_client_secret"`
	Scopes                []string      `json:"scopes"`
	Enabled               bool          `json:"enabled"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}
