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
	BrandSlug        string          `json:"brand_slug" validate:"required,min=3,max=32"`
	Name             string          `json:"name" validate:"required,min=3"`
	Domain           *string         `json:"domain" validate:"omitempty,url"`
	DomainVerifiedAt *time.Time      `json:"domain_verified_at"`
	Metadata         json.RawMessage `json:"metadata"`
	CreatedAt        time.Time       `json:"created_at"`
	DeletedAt        *time.Time      `json:"deleted_at"`
}

func NewProject(ownerID uuid.UUID, slug, name string, domain *string, orgID *uuid.UUID) (*Project, error) {
	p := &Project{
		OrganizationID:   orgID,
		OwnerID:          ownerID,
		Name:             name,
		BrandSlug:        slug,
		Domain:           domain,
		Metadata:         json.RawMessage("{}"),
		DomainVerifiedAt: nil,
	}
	return p, validate.Struct(p)
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

type ProjectMember struct {
	ProjectID uuid.UUID       `json:"project_id"`
	ActorID   uuid.UUID       `json:"actor_id"`
	Role      ProjectRole     `json:"role"`
	Metadata  json.RawMessage `json:"metadata"`
	JoinedAt  time.Time       `json:"joined_at"`
}

func NewProjectMember(projectID, actorID uuid.UUID, role ProjectRole) (*ProjectMember, error) {
	pm := &ProjectMember{
		ProjectID: projectID,
		ActorID:   actorID,
		Role:      role,
		Metadata:  json.RawMessage("{}"),
	}
	return pm, validate.Struct(pm)
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

type CreateProjectRequest struct {
	Name      string  `json:"name"             validate:"required,min=3"`
	Domain    *string `json:"domain,omitempty" validate:"omitempty,url"`
	BrandSlug string  `json:"brand_slug"`
}

func (r CreateProjectRequest) ToInput(orgID *uuid.UUID) CreateProjectInput {
	return CreateProjectInput{
		OrganizationID: orgID,
		Name:           r.Name,
		Domain:         r.Domain,
		BrandSlug:      r.BrandSlug,
	}
}

type CreateProjectInput struct {
	OrganizationID *uuid.UUID
	Name           string
	Domain         *string
	BrandSlug      string
}

type AddProjectMemberRequest struct {
	ActorEmail string      `json:"actor_email"`
	Role       ProjectRole `json:"role"`
}

func (r *AddProjectMemberRequest) ToInput(projectID uuid.UUID) AddProjectMemberInput {
	return AddProjectMemberInput{
		ActorEmail: r.ActorEmail,
		Role:       r.Role,
		ProjectID:  projectID,
	}
}

type AddProjectMemberInput struct {
	ActorEmail string      `json:"actor_email"`
	Role       ProjectRole `json:"role"`
	ProjectID  uuid.UUID   `json:"project_id"`
}

type RemoveProjectMemberRequest struct {
	ActorEmail string `json:"actor_email"`
}

func (r *RemoveProjectMemberRequest) ToInput(projectID uuid.UUID) RemoveProjectMemberInput {
	return RemoveProjectMemberInput{
		ActorEmail: r.ActorEmail,
		ProjectID:  projectID,
	}
}

type RemoveProjectMemberInput struct {
	ActorEmail string    `json:"actor_email"`
	ProjectID  uuid.UUID `json:"project_id"`
}
