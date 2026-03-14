package inbounds

import (
	"GoAuth/internal/domain/project"
	"GoAuth/internal/domain/project_users"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectServiceInput struct {
	ProjectID   uuid.UUID
	ProjectName string
	Domain      string
	Metadata    json.RawMessage
}

type OutputProject struct {
	ID          uuid.UUID
	ProjectName string
	OwnerID     uuid.UUID
	Metadata    json.RawMessage
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func OutputProjectSliceFromProjectSlice(src []project.Project) []OutputProject {
	dst := make([]OutputProject, 0, len(src))
	for _, p := range src {
		dst = append(dst, *OutputProjectFromProject(&p))
	}
	return dst
}

func OutputProjectFromProject(p *project.Project) *OutputProject {
	return &OutputProject{
		ID:          p.ID,
		ProjectName: p.ProjectName,
		OwnerID:     p.OwnerID,
		Metadata:    p.Metadata,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type OutputProjectUser struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Email       string
	UserType    string
	Metadata    *json.RawMessage
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt *time.Time
	IsVerified  bool
	VerifiedAt  *time.Time
}

func OutputProjectUserSliceFromProjectUserSlice(src []project_users.ProjectUser) []OutputProjectUser {
	dst := make([]OutputProjectUser, 0, len(src))
	for _, u := range src {
		dst = append(dst, *OutputProjectUserFromProjectUser(&u))
	}
	return dst
}

func OutputProjectUserFromProjectUser(u *project_users.ProjectUser) *OutputProjectUser {
	if u == nil {
		return nil
	}
	return &OutputProjectUser{
		ID:          u.ID,
		ProjectID:   u.ProjectID,
		Email:       u.Email,
		UserType:    u.UserType,
		Metadata:    u.Metadata,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
		IsVerified:  u.IsVerified,
		VerifiedAt:  u.VerifiedAt,
	}
}
