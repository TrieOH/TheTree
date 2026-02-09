package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"
)

type ProjectResponse struct {
	ID          string          `json:"id"`
	ProjectName string          `json:"project_name"`
	OwnerID     string          `json:"owner_id"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ProjectUserResponse struct {
	ID          string          `json:"id"`
	ProjectID   string          `json:"project_id"`
	Email       string          `json:"email"`
	UserType    string          `json:"user_type"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastLoginAt *time.Time      `json:"last_login_at,omitempty"`
	IsVerified  bool            `json:"is_verified"`
	VerifiedAt  *time.Time      `json:"verified_at,omitempty"`
}

func ProjectUserSliceToProjectUserResponseSlice(src []inbounds.OutputProjectUser) []ProjectUserResponse {
	dst := make([]ProjectUserResponse, 0, len(src))
	for _, u := range src {
		dst = append(dst, ProjectUserToResponse(&u))
	}
	return dst
}

func ProjectUserToResponse(u *inbounds.OutputProjectUser) ProjectUserResponse {
	if u == nil {
		return ProjectUserResponse{}
	}
	var meta json.RawMessage
	if u.Metadata != nil {
		meta = *u.Metadata
	}
	return ProjectUserResponse{
		ID:          u.ID.String(),
		ProjectID:   u.ProjectID.String(),
		Email:       u.Email,
		UserType:    u.UserType,
		Metadata:    meta,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
		IsVerified:  u.IsVerified,
		VerifiedAt:  u.VerifiedAt,
	}
}

func ProjectSliceToProjectResponseSlice(src []inbounds.OutputProject) []ProjectResponse {
	dst := make([]ProjectResponse, 0, len(src))
	for _, p := range src {
		dst = append(dst, ProjectToResponse(&p))
	}
	return dst
}

func ProjectToResponse(r *inbounds.OutputProject) ProjectResponse {
	if r == nil {
		return ProjectResponse{}
	}
	return ProjectResponse{
		ID:          r.ID.String(),
		ProjectName: r.ProjectName,
		OwnerID:     r.OwnerID.String(),
		Metadata:    r.Metadata,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
