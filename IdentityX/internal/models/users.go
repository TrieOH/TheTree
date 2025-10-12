package models

import (
	"github.com/google/uuid"
	"time"
)

type Users struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserSubJWT struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UserTokens struct {
	AccessTokenString  string `json:"access_token"`
	RefreshTokenString string `json:"refresh_token"`
}
