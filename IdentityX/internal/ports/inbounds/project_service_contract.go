package inbounds

import (
	"GoAuth/internal/domain/project"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectServiceInput struct {
	ProjectID   uuid.UUID
	ProjectName string
	Metadata    json.RawMessage
}

type OutputProject struct {
	ID          uuid.UUID
	ProjectName string
	OwnerID     uuid.UUID
	Metadata    json.RawMessage
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func OutputProjectSliceFromProjectSlice(src []project.Project) []OutputProject {
	dst := make([]OutputProject, 0, len(src))
	for _, p := range src {
		dst = append(dst, *OutputProjectFromProject(&p))
	}
	return dst
}

func OutputProjectFromProject(p *project.Project) *OutputProject {
	return &OutputProject{
		ID:          p.ID,
		ProjectName: p.ProjectName,
		OwnerID:     p.OwnerID,
		Metadata:    p.Metadata,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type ErrGeneratingProjectKeys struct {
	Cause error
}

func (e ErrGeneratingProjectKeys) Error() string {
	return "error generating project keys"
}

type ErrParsingProjectPublicKey struct {
	Cause error
}

func (e ErrParsingProjectPublicKey) Error() string {
	return "error parsing project public key"
}

type ErrNotProjectOwner struct {
	Msg string
}

func (e ErrNotProjectOwner) Error() string {
	return e.Msg
}

type ErrFlowIDIsReserved struct {
	Reserved string
}

func (e ErrFlowIDIsReserved) Error() string {
	return "flow id can't be the reserved keyword '" + e.Reserved + "'"
}

type ErrFlowIDSchemaTypeConflict struct{}

func (e ErrFlowIDSchemaTypeConflict) Error() string {
	return "schema with this flow ID already exists in this type"
}

type ErrSchemaNotOwned struct {
	Msg string
}

func (e ErrSchemaNotOwned) Error() string {
	return e.Msg
}

type ErrPublishSchemaPublished struct{}

func (e ErrPublishSchemaPublished) Error() string {
	return "cannot publish a schema that isn't a draft"
}

type ErrPublishSchemaArchived struct{}

func (e ErrPublishSchemaArchived) Error() string {
	return "cannot publish a schema that isn't a draft"
}

type ErrSchemaInvalidStatus struct {
	Status string
}

func (e ErrSchemaInvalidStatus) Error() string {
	return "CATASTROPHIC: schema found with no valid status: " + e.Status
}

type ErrSchemaNoPublishedVersions struct {
	Msg string
}

func (e ErrSchemaNoPublishedVersions) Error() string {
	return e.Msg
}

type ErrSchemaOnlyDraft struct {
	Msg string
}

func (e ErrSchemaOnlyDraft) Error() string {
	return e.Msg
}

type ErrSchemaOnlyArchived struct {
	Msg string
}

func (e ErrSchemaOnlyArchived) Error() string {
	return e.Msg
}
