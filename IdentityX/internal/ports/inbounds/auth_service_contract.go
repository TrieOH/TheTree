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

	IsUpToDate bool
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

type FormResponse struct {
	SchemaID      uuid.UUID
	Title         string
	FlowID        string
	SchemaType    string
	VersionID     uuid.UUID
	VersionNumber int
	Fields        []FormField
}
