package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID        `json:"id"`
	OwnerID   uuid.UUID        `json:"owner_id"   validate:"required"`
	Name      string           `json:"name"       validate:"required,min=3"`
	Slug      string           `json:"slug"       validate:"required,min=2"`
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
	ActorID        uuid.UUID        `json:"actor_id"`
	Role           OrganizationRole `json:"role"`
	Metadata       *json.RawMessage `json:"metadata"`
	JoinedAt       time.Time        `json:"joined_at"`
}

// TODO: kill this constructor — build the struct directly and validate at use sites.
func NewOrganization(ownerID uuid.UUID, name, slug string) (*Organization, error) {
	f := &Organization{
		OwnerID: ownerID,
		Name:    name,
		Slug:    slug,
	}
	err := validate.Struct(f)
	if err != nil {
		return nil, err
	}
	return f, nil
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
	Name      string  `json:"name"             validate:"required,min=3"`
	Domain    *string `json:"domain,omitempty" validate:"omitempty,url"`
	BrandSlug string  `json:"brand_slug"`
}

func (r CreateOrgProjectRequest) ToInput(orgID uuid.UUID) CreateOrgProjectInput {
	return CreateOrgProjectInput{
		OrganizationID: orgID,
		Name:           r.Name,
		Domain:         r.Domain,
		BrandSlug:      r.BrandSlug,
	}
}

type CreateOrgProjectInput struct {
	OrganizationID uuid.UUID
	Name           string
	Domain         *string
	BrandSlug      string `json:"brand_slug"`
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
