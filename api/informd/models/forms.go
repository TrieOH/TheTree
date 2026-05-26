package models

import (
	"time"

	"github.com/google/uuid"
)

type FormStatus string

const (
	FormStatusDraft    FormStatus = "draft"
	FormStatusOpen     FormStatus = "open"
	FormStatusClosed   FormStatus = "closed"
	FormStatusArchived FormStatus = "archived"
)

type Form struct {
	ID          uuid.UUID  `json:"id"`
	NamespaceID *uuid.UUID `json:"namespace_id"`
	OwnerID     uuid.UUID  `json:"owner_id" validate:"required"`
	CreatedBy   uuid.UUID  `json:"created_by" validate:"required"`
	Title       string     `json:"title"    validate:"required"`
	Status      FormStatus `json:"status"`
	OpenedAt    *time.Time `json:"opened_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewForm(namespaceID *uuid.UUID, ownerID, createdBy uuid.UUID, title string) (*Form, error) {
	f := &Form{
		NamespaceID: namespaceID,
		OwnerID:     ownerID,
		CreatedBy:   createdBy,
		Title:       title,
		Status:      FormStatusDraft,
	}
	if err := validate.Struct(f); err != nil {
		return nil, err
	}
	return f, nil
}

type CreateFormRequest struct {
	Title string `json:"title" validate:"required"`
}

type FormMemberRole string

const (
	FormMemberRoleViewer FormMemberRole = "viewer"
	FormMemberRoleEditor FormMemberRole = "editor"
	FormMemberRoleAdmin  FormMemberRole = "admin"
	FormMemberRoleOwner  FormMemberRole = "owner"
)

type FormMember struct {
	UserID  uuid.UUID      `json:"user_id"`
	FormID  uuid.UUID      `json:"form_id"`
	Role    FormMemberRole `json:"role"`
	AddedAt time.Time      `json:"added_at"`
	AddedBy uuid.UUID      `json:"added_by"`
}

type AddFormMemberRequest struct {
	UserID uuid.UUID      `json:"user_id"`
	Role   FormMemberRole `json:"role"`
}

func (r AddFormMemberRequest) ToInput(formID uuid.UUID) AddFormMemberInput {
	return AddFormMemberInput{
		UserID: r.UserID,
		FormID: formID,
		Role:   r.Role,
	}
}

func (r AddFormMemberRequest) ToNamespaceInput(namespaceID, formID uuid.UUID) AddNamespaceFormMemberInput {
	return AddNamespaceFormMemberInput{
		UserID:      r.UserID,
		NamespaceID: namespaceID,
		FormID:      formID,
		Role:        r.Role,
	}
}

type AddFormMemberInput struct {
	UserID uuid.UUID      `json:"user_id"`
	FormID uuid.UUID      `json:"form_id"`
	Role   FormMemberRole `json:"role"`
}

type AddNamespaceFormMemberInput struct {
	UserID      uuid.UUID      `json:"user_id"`
	NamespaceID uuid.UUID      `json:"namespace_id"`
	FormID      uuid.UUID      `json:"form_id"`
	Role        FormMemberRole `json:"role"`
}

type RemoveFormMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

func (r RemoveFormMemberRequest) ToInput(formID uuid.UUID) RemoveFormMemberInput {
	return RemoveFormMemberInput{
		UserID: r.UserID,
		FormID: formID,
	}
}

func (r RemoveFormMemberRequest) ToNamespaceInput(namespaceID, formID uuid.UUID) RemoveNamespaceFormMemberInput {
	return RemoveNamespaceFormMemberInput{
		UserID:      r.UserID,
		NamespaceID: namespaceID,
		FormID:      formID,
	}
}

type RemoveFormMemberInput struct {
	UserID uuid.UUID `json:"user_id"`
	FormID uuid.UUID `json:"form_id"`
}

type RemoveNamespaceFormMemberInput struct {
	UserID      uuid.UUID `json:"user_id"`
	NamespaceID uuid.UUID `json:"namespace_id"`
	FormID      uuid.UUID `json:"form_id"`
}
