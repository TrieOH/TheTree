package outbounds

import (
	"github.com/google/uuid"
)

type VerificationEmailData struct {
	UserID uuid.UUID
	Email  string
	Token  string
	Locale string
}

type PasswordResetEmailData struct {
	UserID string
	Email  string
	Token  string
	Locale string
}

type EmailRenderer interface {
	Verification(data VerificationEmailData) (Email, error)
	PasswordReset(data PasswordResetEmailData) (Email, error)
}
