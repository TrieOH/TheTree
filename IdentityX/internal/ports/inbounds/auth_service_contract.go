package inbounds

import (
	"encoding/json"
	"net/http"
	"time"
)

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
