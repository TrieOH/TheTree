package models

import "github.com/google/uuid"

type CredentialType string

const (
	TokenCredentialType  CredentialType = "token"
	ApiKeyCredentialType CredentialType = "api_key"
)

type IDXRegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

type IDXRegisterInput struct {
	Email     string
	Password  string
	ProjectID *uuid.UUID
}

func (r IDXRegisterRequest) ToInput(projectID *uuid.UUID) IDXRegisterInput {
	return IDXRegisterInput{
		Email:     r.Email,
		Password:  r.Password,
		ProjectID: projectID,
	}
}

type IDXLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd,min=8"`
}

type IDXLoginInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ProjectID *uuid.UUID
}

func (r IDXLoginRequest) ToInput(projectID *uuid.UUID) IDXLoginInput {
	return IDXLoginInput{
		Email:     r.Email,
		Password:  r.Password,
		ProjectID: projectID,
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

type ProjectIDQueryParam struct {
	ProjectID *uuid.UUID `fun_query:"project_id"`
}
