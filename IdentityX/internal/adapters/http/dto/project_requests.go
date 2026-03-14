package dto

import (
	"encoding/json"
)

type CreateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"required,max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

type UpdateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}
