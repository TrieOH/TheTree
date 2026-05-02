package contracts

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSub struct {
	ID         uuid.UUID        `json:"id"`
	Email      string           `json:"email"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	UserType   string           `json:"user_type"`
	Metadata   *json.RawMessage `json:"metadata"`
	SessionID  uuid.UUID        `json:"session_id"`
	UserAgent  string           `json:"user_agent"`
	UserIP     string           `json:"user_ip"`
	IsVerified bool             `json:"is_verified"`
	FamilyID   uuid.UUID        `json:"family_id"`
	VerifiedAt *time.Time       `json:"verified_at"`
}

type AccessClaims struct {
	Sub AccessSub `json:"sub"`
	jwt.RegisteredClaims
}

type RefreshSub struct {
	AccessJTI uuid.UUID `json:"access_jti"`
	FamilyID  uuid.UUID `json:"family_id"`
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

type RegisterInput struct {
	Email     string
	Password  string
	ProjectID *uuid.UUID // nil = client
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

func (r RegisterUserRequest) ToInput(projectID *uuid.UUID) RegisterInput {
	return RegisterInput{
		Email:     r.Email,
		Password:  r.Password,
		ProjectID: projectID,
	}
}

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	Agent     string
	ProjectID *uuid.UUID // nil = client
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

func (r LoginUserRequest) ToInput(projectID *uuid.UUID, agent, ip string) LoginInput {
	return LoginInput{
		Email:     r.Email,
		Password:  r.Password,
		ProjectID: projectID,
		Agent:     agent,
		IP:        ip,
	}
}

type UserTokensOutput struct {
	AccessTokenString  string
	RefreshTokenString string
	AccessExpiresAt    time.Time
	RefreshExpiresAt   time.Time
	Domain             string
}

type UserTokensResponse struct {
	AccessTokenString  string
	RefreshTokenString string
	AccessExpiresAt    time.Time
	RefreshExpiresAt   time.Time
	Domain             string
}

func (r UserTokensOutput) ToResponse() UserTokensResponse {
	return UserTokensResponse{
		AccessTokenString:  r.AccessTokenString,
		RefreshTokenString: r.RefreshTokenString,
		AccessExpiresAt:    r.AccessExpiresAt,
		RefreshExpiresAt:   r.RefreshExpiresAt,
		Domain:             r.Domain,
	}
}

type RefreshInput struct {
	RefreshCookie string
	Agent         string
	IP            string
}

func ToRefreshInput(refreshCookie, Agent, IP string) RefreshInput {
	return RefreshInput{
		RefreshCookie: refreshCookie,
		Agent:         Agent,
		IP:            IP,
	}
}
