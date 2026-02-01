package inbounds

import (
	"GoAuth/internal/domain/version"
	"time"

	"github.com/google/uuid"
)

type SchemaVersionServiceInput struct {
	SchemaID      uuid.UUID
	ProjectID     uuid.UUID
	VersionID     *uuid.UUID
	VersionNumber int
}

type VersionVerboseOutput struct {
	SchemaVersionOutput
	Fields []OutputField
}

type SchemaVersionOutput struct {
	ID               uuid.UUID
	SchemaID         uuid.UUID
	BasedOnVersionID *uuid.UUID
	VersionNumber    int
	Status           version.Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func SchemaVersionToOutput(out *version.Version) *SchemaVersionOutput {
	if out == nil {
		return nil
	}
	return &SchemaVersionOutput{
		ID:               out.ID,
		SchemaID:         out.SchemaID,
		BasedOnVersionID: out.BasedOnVersionID,
		VersionNumber:    out.VersionNumber,
		Status:           out.Status,
		CreatedAt:        out.CreatedAt,
		UpdatedAt:        out.UpdatedAt,
	}
}

type ErrDraftVersionOnNonPublished struct{}

func (e ErrDraftVersionOnNonPublished) Error() string {
	return "new versions can only be drafted from published versions"
}

type ErrPublishSchemaNonExistentVersionDraft struct{}

func (e ErrPublishSchemaNonExistentVersionDraft) Error() string {
	return "cannot publish a schema with a version draft that doesn't exist"
}

type ErrPublishVersionPublished struct{}

func (e ErrPublishVersionPublished) Error() string {
	return "cannot publish a schema version that is already published"
}

type ErrPublishVersionArchived struct{}

func (e ErrPublishVersionArchived) Error() string {
	return "cannot publish a schema version that is archived"
}

type ErrPublishVersionNotDraft struct{}

func (e ErrPublishVersionNotDraft) Error() string {
	return "cannot publish a schema version that isn't a draft"
}

type ErrPublishNonExistentVersion struct{}

func (e ErrPublishNonExistentVersion) Error() string {
	return "cannot publish a non-existent schema version"
}

type ErrPublishVersionInvalidStatus struct{}

func (e ErrPublishVersionInvalidStatus) Error() string {
	return "CATASTROPHIC: schema version found with no valid status"
}

type ErrPublishVersionNoChanges struct{}

func (e ErrPublishVersionNoChanges) Error() string {
	return "cannot publish a version with no changes"
}
