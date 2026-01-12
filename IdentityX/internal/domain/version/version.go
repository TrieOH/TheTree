package version

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Version struct {
	ID               uuid.UUID
	SchemaID         uuid.UUID
	VersionNumber    int
	Status           Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
	BasedOnVersionID *uuid.UUID
}

func (v Version) CanRegister() error {
	if v.Status == StatusDraft {
		return ErrRegisterOnVersionDraft{}
	}
	if v.Status == StatusArchived {
		return ErrRegisterOnVersionArchive{}
	}
	return nil
}
