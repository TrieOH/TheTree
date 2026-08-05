package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type EmailTemplateKind string

const (
	VerifyEmailTemplateKind EmailTemplateKind = "verify"
	ResetEmailTemplateKind  EmailTemplateKind = "reset"
)

// AllEmailTemplateKinds is the set of template kinds a project can
// override. The effective template for a project is its override (if any)
// or the baked-in default.
var AllEmailTemplateKinds = []EmailTemplateKind{VerifyEmailTemplateKind, ResetEmailTemplateKind}

type EmailTemplate struct {
	ID        uuid.UUID         `json:"id"`
	ProjectID uuid.UUID         `json:"project_id"`
	Kind      EmailTemplateKind `json:"kind"`
	Subject   string            `json:"subject"`
	Body      string            `json:"body"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// EffectiveEmailTemplate is a template as served to clients: the project
// override when one exists, otherwise the built-in default, with Source
// telling the caller which one it is.
type EffectiveEmailTemplate struct {
	Kind    EmailTemplateKind `json:"kind"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	Source  string            `json:"source"` // "default" | "override"
}

type ActionTokenPurpose string

const (
	EmailVerifyActionTokenPurpose   ActionTokenPurpose = "email_verify"
	PasswordResetActionTokenPurpose ActionTokenPurpose = "password_reset"
)

type ActionToken struct {
	JTI       uuid.UUID          `json:"jti"`
	Purpose   ActionTokenPurpose `json:"purpose"`
	ActorID   uuid.UUID          `json:"actor_id"`
	ExpiresAt time.Time          `json:"expires_at"`
	UsedAt    *time.Time         `json:"used_at"`
	CreatedAt time.Time          `json:"created_at"`
}

// ActionTokenClaims is the JWT carried by verify/reset links. Subject is
// the actor ID, ID the jti (both uuid strings); Purpose scopes the token to
// exactly one flow; ProjectID distinguishes project actors from platform
// actors (it also travels on the URL for the SPA, but the claim is the
// source of truth at redemption).
type ActionTokenClaims struct {
	jwt.RegisteredClaims

	Purpose   string     `json:"purpose"`
	ProjectID *uuid.UUID `json:"project_id,omitempty"`
}
