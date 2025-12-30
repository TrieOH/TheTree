package dto

import "encoding/json"

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
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
