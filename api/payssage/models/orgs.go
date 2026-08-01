package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID  `json:"id"`
	OwnerID   uuid.UUID  `json:"owner_id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type OrganizationRole string

const (
	OrganizationRoleMember OrganizationRole = "member"
	OrganizationRoleAdmin  OrganizationRole = "admin"
	OrganizationRoleOwner  OrganizationRole = "owner"
)

func (r OrganizationRole) Rank() int {
	switch r {
	case OrganizationRoleMember:
		return 0
	case OrganizationRoleAdmin:
		return 1
	case OrganizationRoleOwner:
		return 2
	default:
		return 0
	}
}

func (r OrganizationRole) String() string { return string(r) }

type OrganizationMember struct {
	OrganizationID uuid.UUID        `json:"organization_id"`
	MemberID       uuid.UUID        `json:"member_id"`
	Role           OrganizationRole `json:"role"`
	JoinedAt       time.Time        `json:"joined_at"`
}

type CreateOrganizationInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type AddOrganizationMemberInput struct {
	ActorEmail     string           `json:"actor_email"`
	Role           OrganizationRole `json:"role"`
	OrganizationID uuid.UUID        `json:"organization_id"`
}

type RemoveOrganizationMemberInput struct {
	ActorEmail     string    `json:"actor_email"`
	OrganizationID uuid.UUID `json:"organization_id"`
}
