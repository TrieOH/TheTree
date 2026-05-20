package models

import (
	"github.com/google/uuid"
)

type ForgotPasswordInput struct {
	Email     string     `json:"email"`
	ProjectID *uuid.UUID `json:"project_id"`
}

type ForgotPasswordRequest struct {
	Email     string     `json:"email" validate:"required,email"`
	ProjectID *uuid.UUID `json:"project_id"`
}

func (r ForgotPasswordRequest) ToInput() ForgotPasswordInput {
	return ForgotPasswordInput{
		Email:     r.Email,
		ProjectID: r.ProjectID,
	}
}

type ResetPasswordInput struct {
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,passwd,min=8,max=72"`
}

func (r ResetPasswordRequest) ToInput(token string) ResetPasswordInput {
	return ResetPasswordInput{
		NewPassword: r.NewPassword,
		Token:       token,
	}
}
