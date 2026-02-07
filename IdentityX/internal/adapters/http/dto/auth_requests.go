package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

type RegisterProjectUserRequest struct {
	Email        string           `json:"email" validate:"required,email,max=255"`
	Password     string           `json:"password" validate:"required,passwd,min=8,max=72"`
	CustomFields *json.RawMessage `json:"custom_fields"`
}

type LoginProjectUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

type UpdateMetadataRequest struct {
	CustomFields *json.RawMessage `json:"custom_fields" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email     string     `json:"email" validate:"required,email"`
	ProjectID *uuid.UUID `json:"project_id"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,passwd,min=8,max=72"`
}
