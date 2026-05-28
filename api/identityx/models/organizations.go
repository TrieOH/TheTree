package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID        `json:"id"`
	OwnerID   uuid.UUID        `json:"owner_id" validate:"required"`
	Name      string           `json:"name" validate:"required,min=3"`
	Slug      string           `json:"slug" validate:"required,min=2"`
	Metadata  *json.RawMessage `json:"metadata"`
	CreatedAt time.Time        `json:"created_at"`
	DeletedAt *time.Time       `json:"deleted_at"`
}

type OrganizationRole string

const (
	OrganizationRoleMember OrganizationRole = "member"
	OrganizationRoleAdmin  OrganizationRole = "admin"
	OrganizationRoleOwner  OrganizationRole = "owner"
)

type OrganizationMember struct {
	OrganizationID uuid.UUID        `json:"organization_id"`
	ActorID        uuid.UUID        `json:"actor_id"`
	Role           OrganizationRole `json:"role"`
	Metadata       *json.RawMessage `json:"metadata"`
	JoinedAt       time.Time        `json:"joined_at"`
}

func NewOrganization(ownerID uuid.UUID, name, slug string) (*Organization, error) {
	f := &Organization{
		OwnerID: ownerID,
		Name:    name,
		Slug:    slug,
	}
	return f, validate.Struct(f)
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
	ActorID uuid.UUID        `json:"actor_id"`
	Role    OrganizationRole `json:"role"`
}

func (r *AddOrganizationMemberRequest) ToInput(orgID uuid.UUID) AddOrganizationMemberInput {
	return AddOrganizationMemberInput{
		ActorID:        r.ActorID,
		Role:           r.Role,
		OrganizationID: orgID,
	}
}

type AddOrganizationMemberInput struct {
	ActorID        uuid.UUID        `json:"actor_id"`
	Role           OrganizationRole `json:"role"`
	OrganizationID uuid.UUID        `json:"organization_id"`
}

type RemoveOrganizationMemberRequest struct {
	ActorID uuid.UUID `json:"actor_id"`
}

func (r *RemoveOrganizationMemberRequest) ToInput(orgID uuid.UUID) RemoveOrganizationMemberInput {
	return RemoveOrganizationMemberInput{
		ActorID:        r.ActorID,
		OrganizationID: orgID,
	}
}

type RemoveOrganizationMemberInput struct {
	ActorID        uuid.UUID `json:"actor_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
}
