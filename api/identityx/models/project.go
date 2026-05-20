package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID       `json:"id"`
	ProjectName string          `json:"project_name"`
	Domain      string          `json:"domain"`
	OwnerID     uuid.UUID       `json:"owner_id"`
	Metadata    json.RawMessage `json:"metadata"`
	IsActive    bool            `json:"is_active"`
	PubKey      string          `json:"pub_key"`
	PrivKey     []byte          `json:"-"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateProjectInput struct {
	ProjectName string          `json:"project_name"`
	Domain      string          `json:"domain"`
	Metadata    json.RawMessage `json:"metadata"`
}

type CreateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"required,max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

func (r CreateProjectRequest) ToInput() CreateProjectInput {
	return CreateProjectInput{
		ProjectName: r.ProjectName,
		Domain:      r.Domain,
		Metadata:    r.Metadata,
	}
}

type UpdateProjectInput struct {
	ProjectID   uuid.UUID       `json:"project_id"`
	ProjectName string          `json:"project_name"`
	Domain      string          `json:"domain"`
	Metadata    json.RawMessage `json:"metadata"`
}

type UpdateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

func (r UpdateProjectRequest) ToInput(projectID uuid.UUID) UpdateProjectInput {
	return UpdateProjectInput{
		ProjectName: r.ProjectName,
		Domain:      r.Domain,
		Metadata:    r.Metadata,
		ProjectID:   projectID,
	}
}
