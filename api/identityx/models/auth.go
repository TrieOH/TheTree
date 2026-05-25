package models

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSub struct {
	ID           uuid.UUID       `json:"id"`
	ProjectID    *uuid.UUID      `json:"project_id"`
	Email        *string         `json:"email"`
	Type         ActorType       `json:"type"`
	Capabilities json.RawMessage `json:"capabilities"`
	Metadata     json.RawMessage `json:"metadata"`
}

type AccessClaims struct {
	Sub AccessSub `json:"sub"`
	jwt.RegisteredClaims
}

type RefreshSub struct {
	AccessJTI uuid.UUID `json:"access_jti"`
}

type RefreshClaims struct {
	Sub RefreshSub `json:"sub"`
	jwt.RegisteredClaims
}

type VerificationSub struct {
	Subject uuid.UUID `json:"subject"`
}

type VerificationClaims struct {
	Sub VerificationSub `json:"sub"`
	jwt.RegisteredClaims
}

type ResetPasswordSub struct {
	Subject   uuid.UUID  `json:"subject"`
	ProjectID *uuid.UUID `json:"project_id"`
}

type ResetPasswordClaims struct {
	Sub ResetPasswordSub `json:"sub"`
	jwt.RegisteredClaims
}

type IDXRegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

type IDXRegisterInput struct {
	Email    string
	Password string
}

func (r IDXRegisterRequest) ToInput() IDXRegisterInput {
	return IDXRegisterInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type IDXLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd,min=8"`
}

type IDXLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r IDXLoginRequest) ToInput() IDXLoginInput {
	return IDXLoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type SetupInput struct {
	Email    string
	Password string
}

func (r IDXLoginRequest) ToSetupInput() SetupInput {
	return SetupInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type UserTokensOutput struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	Domain           string    `json:"domain"`
}

type UserTokensResponse struct {
	AccessTokenString  string    `json:"access_token_string"`
	RefreshTokenString string    `json:"refresh_token_string"`
	AccessExpiresAt    time.Time `json:"access_expires_at"`
	RefreshExpiresAt   time.Time `json:"refresh_expires_at"`
	Domain             string    `json:"domain"`
}

func (r UserTokensOutput) ToResponse() UserTokensResponse {
	return UserTokensResponse{
		AccessTokenString:  r.AccessToken,
		RefreshTokenString: r.RefreshToken,
		AccessExpiresAt:    r.AccessExpiresAt,
		RefreshExpiresAt:   r.RefreshExpiresAt,
		Domain:             r.Domain,
	}
}

type RefreshInput struct {
	RefreshCookie string `json:"refresh_cookie"`
	Agent         string `json:"agent"`
	IP            string `json:"ip"`
}

func ToRefreshInput(refreshCookie, Agent, IP string) RefreshInput {
	return RefreshInput{
		RefreshCookie: refreshCookie,
		Agent:         Agent,
		IP:            IP,
	}
}
