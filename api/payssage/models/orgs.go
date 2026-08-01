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

type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (r *CreateOrganizationRequest) ToInput() CreateOrganizationInput {
	return CreateOrganizationInput{
		Name: r.Name,
		Slug: r.Slug,
	}
}

type CreateOrganizationInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type AddOrganizationMemberRequest struct {
	ActorEmail string           `json:"actor_email"`
	Role       OrganizationRole `json:"role"`
}

func (r *AddOrganizationMemberRequest) ToInput(orgID uuid.UUID) AddOrganizationMemberInput {
	return AddOrganizationMemberInput{
		ActorEmail:     r.ActorEmail,
		Role:           r.Role,
		OrganizationID: orgID,
	}
}

type AddOrganizationMemberInput struct {
	ActorEmail     string           `json:"actor_email"`
	Role           OrganizationRole `json:"role"`
	OrganizationID uuid.UUID        `json:"organization_id"`
}

type RemoveOrganizationMemberRequest struct {
	ActorEmail string `json:"actor_email"`
}

func (r *RemoveOrganizationMemberRequest) ToInput(orgID uuid.UUID) RemoveOrganizationMemberInput {
	return RemoveOrganizationMemberInput{
		ActorEmail:     r.ActorEmail,
		OrganizationID: orgID,
	}
}

type RemoveOrganizationMemberInput struct {
	ActorEmail     string    `json:"actor_email"`
	OrganizationID uuid.UUID `json:"organization_id"`
}
