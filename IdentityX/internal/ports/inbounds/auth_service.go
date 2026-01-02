package inbounds

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, in RegisterUserInput) error
	Login(ctx context.Context, in LoginUserInput) (*UserTokensOutput, error)
	Logout(ctx context.Context) error
	Refresh(ctx context.Context, in RefreshInput) (*UserTokensOutput, error)
	RegisterProjectUser(ctx context.Context, in ProjectRegisterInput) error
	LoginProjectUser(ctx context.Context, in ProjectLoginInput) (*UserTokensOutput, error)
}

type RegisterUserInput struct {
	Email    string
	Password string
}

type LoginUserInput struct {
	Email    string
	Password string

	Agent string
	IP    string
}

type UserTokensOutput struct {
	AccessTokenString  string
	RefreshTokenString string

	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type ProjectRegisterInput struct {
	Email        string
	Password     string
	CustomFields json.RawMessage
	ProjectID    string
}

type ProjectLoginInput struct {
	Email     string
	Password  string
	ProjectID string
	IP        string
	Agent     string
}

type RefreshInput struct {
	RefreshCookie *http.Cookie
	Agent         string
	IP            string
}
