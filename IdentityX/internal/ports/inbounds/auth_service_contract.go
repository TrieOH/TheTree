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

type ErrHashingPassword struct {
	Cause error
}

func (e ErrHashingPassword) Error() string {
	return "error hashing user password"
}

type ErrGeneratingUUID struct {
	Cause    error
	Location string
}

func (e ErrGeneratingUUID) Error() string {
	return "error generating UUID V7 at " + e.Location
}

type ErrEmptyFlowID struct{}

func (e ErrEmptyFlowID) Error() string {
	return "flow id can't be empty"
}

type ErrEmptySchemaType struct{}

func (e ErrEmptySchemaType) Error() string {
	return "schema type can't be empty"
}

type ErrInvalidSchemaType struct{}

func (e ErrInvalidSchemaType) Error() string {
	return "invalid schema type"
}

type ErrInvalidFlowID struct {
	Why string
}

func (e ErrInvalidFlowID) Error() string {
	if e.Why == "" {
		return "invalid flow ID"
	}
	return "invalid flow ID: " + e.Why
}

type ErrCustomFieldsNotAllowed struct{}

func (e ErrCustomFieldsNotAllowed) Error() string {
	return "custom fields are not allowed for core schema"
}

type ErrTokenReuseNotAllowed struct {
	TokenType string
}

func (e ErrTokenReuseNotAllowed) Error() string {
	return e.TokenType + " token reuse not allowed"
}

type ErrTokenUserMismatch struct {
	TokenType string
}

func (e ErrTokenUserMismatch) Error() string {
	return e.TokenType + " token user mismatch"
}

type ErrFailedToRetrieveJWKS struct {
	Cause error
}

func (e ErrFailedToRetrieveJWKS) Error() string {
	return "failed to retrieve JWKS"
}

type ErrSchemaRegisterValidation struct {
	Errors []string
}

func (e ErrSchemaRegisterValidation) Error() string {
	return "error validating fields for schema register"
}
