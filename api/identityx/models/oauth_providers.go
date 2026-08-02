package models

import (
	"time"

	"github.com/google/uuid"
)

// CreateOAuthProviderInput is the payload for configuring a provider for a
// project. The secret is encrypted before it touches the database.
type CreateOAuthProviderInput struct {
	ProjectID    uuid.UUID     `json:"-"`
	Provider     OAuthProvider `json:"provider"      validate:"required,oneof=google github"`
	ClientID     string        `json:"client_id"     validate:"required,min=1"`
	ClientSecret string        `json:"client_secret" validate:"required,min=1"`
}

// UpdateOAuthProviderInput is the partial payload for editing a provider.
// A nil ClientID / ClientSecret leaves the stored value untouched.
type UpdateOAuthProviderInput struct {
	ID           uuid.UUID `json:"-"`
	ClientID     *string   `json:"client_id"     validate:"omitempty,min=1"`
	ClientSecret *string   `json:"client_secret" validate:"omitempty,min=1"`
}

// OAuthLoginState is one connect attempt: the opaque state token handed to
// the provider, which credentials scope (project or platform) it was created
// for, and when it stops being valid. Consumed — deleted — by the callback.
type OAuthLoginState struct {
	ID        uuid.UUID     `json:"id"`
	State     string        `json:"-"`
	Provider  OAuthProvider `json:"provider"`
	ProjectID *uuid.UUID    `json:"project_id"`
	CreatedAt time.Time     `json:"created_at"`
	ExpiresAt time.Time     `json:"expires_at"`
}
