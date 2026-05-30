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

// CreateOrgProjectRequest is the HTTP request body for creating an org-scoped project.
type CreateOrgProjectRequest struct {
	Name   string  `json:"name"             validate:"required,min=3"`
	Domain *string `json:"domain,omitempty" validate:"omitempty,url"`
}

func (r CreateOrgProjectRequest) ToInput(orgID uuid.UUID) CreateOrgProjectInput {
	return CreateOrgProjectInput{
		OrganizationID: orgID,
		Name:           r.Name,
		Domain:         r.Domain,
	}
}

type CreateOrgProjectInput struct {
	OrganizationID uuid.UUID
	Name           string
	Domain         *string
}

// AddOrgProjectMemberRequest is the HTTP request body for adding a member to an org-scoped project.
type AddOrgProjectMemberRequest struct {
	ActorEmail string      `json:"actor_email"`
	Role       ProjectRole `json:"role"`
}

func (r *AddOrgProjectMemberRequest) ToInput(orgID, projectID uuid.UUID) AddOrgProjectMemberInput {
	return AddOrgProjectMemberInput{
		ActorEmail:     r.ActorEmail,
		Role:           r.Role,
		OrganizationID: orgID,
		ProjectID:      projectID,
	}
}

type AddOrgProjectMemberInput struct {
	ActorEmail     string
	Role           ProjectRole
	OrganizationID uuid.UUID
	ProjectID      uuid.UUID
}

// RemoveOrgProjectMemberRequest is the HTTP request body for removing a member from an org-scoped project.
type RemoveOrgProjectMemberRequest struct {
	ActorEmail string `json:"actor_email"`
}

func (r *RemoveOrgProjectMemberRequest) ToInput(orgID, projectID uuid.UUID) RemoveOrgProjectMemberInput {
	return RemoveOrgProjectMemberInput{
		ActorEmail:     r.ActorEmail,
		OrganizationID: orgID,
		ProjectID:      projectID,
	}
}

type RemoveOrgProjectMemberInput struct {
	ActorEmail     string
	OrganizationID uuid.UUID
	ProjectID      uuid.UUID
}
