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

type VersionStatus string

const (
	VersionStatusDraft     VersionStatus = "draft"
	VersionStatusPublished VersionStatus = "published"
	VersionStatusArchived  VersionStatus = "archived"
)

type Version struct {
	ID        uuid.UUID
	SchemaID  uuid.UUID
	Version   int
	Status    VersionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
