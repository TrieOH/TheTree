package contracts

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
	OwnerID     uuid.UUID  `json:"owner_id"          validate:"required"`
	Name        string     `json:"name"              validate:"required"` // FIXME : Change to title
	Status      FormStatus `json:"status"`
	OpenedAt    *time.Time `json:"opened_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewForm(namespaceID *uuid.UUID, ownerID uuid.UUID, title string) (*Form, error) {
	f := &Form{
		NamespaceID: namespaceID,
		OwnerID:     ownerID,
		Name:        title,
		Status:      FormStatusDraft,
	}
	if err := validate.Struct(f); err != nil {
		return nil, err
	}
	return f, nil
}
