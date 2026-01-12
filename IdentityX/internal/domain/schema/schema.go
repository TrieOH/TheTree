package schema

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	Core       Type = "core"
	Context    Type = "context"
	SubContext Type = "sub-context"
)

func IsValidSchemaType(s string) bool {
	switch Type(s) {
	case Core, Context, SubContext:
		return true
	default:
		return false
	}
}

type ReservedFlowID string

const NoFlowID ReservedFlowID = "none"

func IsFlowIDReserved(flowID string) bool {
	switch ReservedFlowID(flowID) {
	case NoFlowID:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Schema struct {
	ID               uuid.UUID
	ProjectID        uuid.UUID
	Title            string
	FlowID           string
	Type             Type
	CurrentVersionID *uuid.UUID
	Status           Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (s Schema) CanRegister() error {
	if s.CurrentVersionID == nil {
		return ErrSchemaNoPublishedVersion{}
	}
	if s.Status == StatusDraft {
		return ErrRegisterOnSchemaDraft{}
	}
	if s.Status == StatusArchived {
		return ErrRegisterOnSchemaArchive{}
	}
	return nil
}

func (s Schema) IsVersion(versionID uuid.UUID) bool {
	if s.CurrentVersionID == nil {
		return false
	}
	if *s.CurrentVersionID != versionID {
		return false
	}
	return true
}
