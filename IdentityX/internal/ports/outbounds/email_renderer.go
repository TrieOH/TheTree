package outbounds

import (
	"context"

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
	Verification(ctx context.Context, data VerificationEmailData) (Email, error)
	PasswordReset(ctx context.Context, data PasswordResetEmailData) (Email, error)
}
