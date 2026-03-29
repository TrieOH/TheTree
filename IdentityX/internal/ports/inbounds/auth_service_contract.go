package inbounds

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	Domain string
}

type ProjectRegisterInput struct {
	Email        string
	Password     string
	CustomFields *json.RawMessage
	ProjectID    uuid.UUID
	SchemaType   string
	FlowID       string
}

type ProjectLoginInput struct {
	Email     string
	Password  string
	ProjectID uuid.UUID
	IP        string
	Agent     string
}

type ProjectLogoutInput struct {
	ProjectID          uuid.UUID
	RefreshTokenCookie *http.Cookie
}

type RefreshInput struct {
	RefreshCookie *http.Cookie
	Agent         string
	IP            string
}

type ErrSchemaRegisterValidation struct {
	Errors []string
}

func (e ErrSchemaRegisterValidation) Error() string {
	return "error validating fields for schema register"
}

type ForgotPasswordInput struct {
	Email     string
	ProjectID *uuid.UUID
}

type ResetPasswordInput struct {
	NewPassword string
	Token       string
}

type ExchangeOutput struct {
	SessionID string    `json:"session_id"`
	TTL       time.Time `json:"ttl"`
}
