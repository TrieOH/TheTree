package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string `json:"-"`
	UserType     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
