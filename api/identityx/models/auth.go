package models

import "github.com/google/uuid"

type CredentialType string

const (
	TokenCredentialType  CredentialType = "token"
	APIKeyCredentialType CredentialType = "api_key"
)

type IDXRegisterInput struct {
	Email     string
	Password  string
	ProjectID *uuid.UUID
}

type IDXLoginInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ProjectID *uuid.UUID
}

type SetupInput struct {
	Email    string
	Password string
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type VerifyEmailInput struct {
	Token string
}

type ResendVerificationInput struct {
	Email     string
	ProjectID *uuid.UUID
}

type ForgotPasswordInput struct {
	Email     string
	ProjectID *uuid.UUID
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}
