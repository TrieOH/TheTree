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

var roleRank = map[OrganizationRole]int{
	OrganizationRoleMember: 1,
	OrganizationRoleAdmin:  2,
	OrganizationRoleOwner:  3,
}

func (r OrganizationRole) AtLeast(min OrganizationRole) bool {
	return roleRank[r] >= roleRank[min]
}

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
