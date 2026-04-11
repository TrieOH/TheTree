package contracts

import (
	"TrieForms/internal/shared/validation"
	"time"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
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
	ProjectID        uuid.UUID  `json:"project_id"`
	OwnerID          uuid.UUID  `json:"owner_id"`
	ScopeID          uuid.UUID  `json:"scope_id"`
	Title            string     `json:"title"`
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

	if err = f.validate(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Form) validate() error {
	return validation.Run(
		validation.RequireUUID("form", "owner_id", f.OwnerID),
		validation.RequireUUID("form", "project_id", f.ProjectID),
		validation.RequireString("form", "title", f.Title),
	)
}

func (f *Form) AddScope(scopeID uuid.UUID) {
	f.ScopeID = scopeID
}
