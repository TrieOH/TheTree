package authz

import (
	"IdentityX/contracts"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type AuthMethod string

const (
	AuthMethodSession AuthMethod = "session"
	AuthMethodApiKey  AuthMethod = "api_key"
)

type Principal struct {
	UserID    uuid.UUID          `json:"user_id"`
	UserType  contracts.UserType `json:"user_type"`
	ProjectID *uuid.UUID         `json:"project_id"`
	SessionID *uuid.UUID         `json:"session_id"`
	Method    AuthMethod         `json:"-"`
}

func NewPrincipal(access *contracts.AccessClaims) (*Principal, error) {
	if access == nil {
		return nil, fun.ErrUnprocessableEntity("missing access claims")
	}
	userType := contracts.UserTypeClient
	if access.Sub.ProjectID != nil {
		userType = contracts.UserTypeProject
	}

	return &Principal{
		UserID:    access.Sub.ID,
		UserType:  userType,
		ProjectID: access.Sub.ProjectID,
		SessionID: &access.Sub.SessionID,
		Method:    AuthMethodSession,
	}, nil
}
