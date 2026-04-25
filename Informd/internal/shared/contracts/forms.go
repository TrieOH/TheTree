package contracts

import (
	"time"

	"github.com/MintzyG/FastUtilitiesNet"
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
	ID               uuid.UUID  `json:"id"`
	ProjectID        uuid.UUID  `json:"project_id"        validate:"required"`
	OwnerID          uuid.UUID  `json:"owner_id"          validate:"required"`
	Title            string     `json:"title"             validate:"required"`
	Status           FormStatus `json:"status"`
	CurrentVersionID *uuid.UUID `json:"current_version_id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	OpenedAt         *time.Time `json:"opened_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	ArchivedAt       *time.Time `json:"archived_at"`
}

func NewForm(projectID, ownerID uuid.UUID, title string) (*Form, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fun.NewErrorf("error generating uuid for form: %s", err.Error()).Internal()
	}

	f := &Form{
		ID:        id,
		ProjectID: projectID,
		OwnerID:   ownerID,
		Title:     title,
		Status:    FormStatusDraft,
	}

	if err = validate.Struct(f); err != nil {
		return nil, err
	}
	return f, nil
}
