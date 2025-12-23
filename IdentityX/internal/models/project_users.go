package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectUser struct {
	ID          uuid.UUID        `json:"id"`
	ProjectID   uuid.UUID        `json:"project_id"`
	Email       string           `json:"email"`
	Password    string           `json:"password"`
	UserType    string           `json:"user_type"`
	Metadata    *json.RawMessage `json:"metadata"`
	IsActive    bool             `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	LastLoginAt time.Time        `json:"last_login_at"`
}

type RegisterProjectUserRequest struct {
	Email        string          `json:"email" validate:"required,email,max=255"`
	Password     string          `json:"password" validate:"required,passwd,min=8,max=72"`
	CustomFields json.RawMessage `json:"custom_fields" validate:"required"`
}

type LoginProjectUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}
