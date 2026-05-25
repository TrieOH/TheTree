package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Organizations struct {
	ID        uuid.UUID       `json:"id"`
	OwnerID   uuid.UUID       `json:"owner_id"`
	Name      string          `json:"name"`
	Slug      string          `json:"slug"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
	DeletedAt time.Time       `json:"deleted_at"`
}

type OrganizationRole string

const (
	OrganizationRoleMember OrganizationRole = "member"
	OrganizationRoleAdmin  OrganizationRole = "admin"
	OrganizationRoleOwner  OrganizationRole = "owner"
)

type OrganizationMembers struct {
	OrganizationID uuid.UUID        `json:"organization_id"`
	ActorID        uuid.UUID        `json:"actor_id"`
	Role           OrganizationRole `json:"role"`
	Metadata       json.RawMessage  `json:"metadata"`
	JoinedAt       time.Time        `json:"joined_at"`
}
