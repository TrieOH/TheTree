package models

import (
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeClient  UserType = "client"
	UserTypeProject UserType = "project"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	UserType     UserType   `json:"user_type"`
	ProjectID    *uuid.UUID `json:"project_id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	IsVerified   bool       `json:"is_verified"`
	VerifiedAt   *time.Time `json:"verified_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
