package version

import (
	"GoAuth/internal/errx"
	"context"
	"time"

	"github.com/MintzyG/fail/v3"
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

func (v Version) CanRegister(ctx context.Context) error {
	if v.Status == StatusDraft {
		return fail.New(errx.ProjectUserRegisterOnSchemaVersionDraft).RecordCtx(ctx)
	}
	if v.Status == StatusArchived {
		return fail.New(errx.ProjectUserRegisterOnSchemaVersionArchived).RecordCtx(ctx)
	}
	return nil
}
