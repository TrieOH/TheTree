package inbounds

import (
	"GoAuth/internal/domain/schema"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SchemaServiceInput struct {
	SchemaType string
	Title      string
	FlowID     string
	ProjectID  uuid.UUID
	SchemaID   uuid.UUID
}

type SchemaVerboseOutput struct {
	SchemaOutput
	Versions []VersionVerboseOutput
}

type SchemaOutput struct {
	ID               uuid.UUID
	ProjectID        uuid.UUID
	Title            string
	FlowID           string
	Type             string
	CurrentVersionID *uuid.UUID
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func SchemaSliceToSchemaOutputSlice(out []schema.Schema) []SchemaOutput {
	if out == nil {
		return nil
	}

	outSlice := make([]SchemaOutput, 0, len(out))
	for _, s := range out {
		outSlice = append(outSlice, SchemaOutput{
			ID:               s.ID,
			ProjectID:        s.ProjectID,
			Title:            s.Title,
			FlowID:           s.FlowID,
			Type:             string(s.Type),
			CurrentVersionID: s.CurrentVersionID,
			Status:           string(s.Status),
			CreatedAt:        s.CreatedAt,
			UpdatedAt:        s.UpdatedAt,
		})
	}
	return outSlice
}

func SchemaToSchemaOutput(out *schema.Schema) *SchemaOutput {
	if out == nil {
		return nil
	}
	return &SchemaOutput{
		ID:               out.ID,
		ProjectID:        out.ProjectID,
		Title:            out.Title,
		FlowID:           out.FlowID,
		Type:             string(out.Type),
		CurrentVersionID: out.CurrentVersionID,
		Status:           string(out.Status),
		CreatedAt:        out.CreatedAt,
		UpdatedAt:        out.UpdatedAt,
	}
}

type FormOutput struct {
	SchemaID      uuid.UUID
	Title         string
	FlowID        string
	SchemaType    string
	VersionID     uuid.UUID
	VersionNumber int
	Status        string
	Fields        []FormField
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FormField struct {
	ID              uuid.UUID
	ObjectID        uuid.UUID
	Key             string
	Type            string
	Owner           string
	Title           string
	Description     *string
	Placeholder     *string
	Required        bool
	Mutable         bool
	DefaultValue    *json.RawMessage
	Position        int
	Options         []FormOption
	VisibilityRules []FormRule
	RequiredRules   []FormRule
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type FormOption struct {
	ID       uuid.UUID
	Value    string
	Label    string
	Position int
}

type FormRule struct {
	ID               uuid.UUID
	DependsOnFieldID uuid.UUID
	Operator         string
	Value            *json.RawMessage
}

type ErrVersionNotPublished struct{}

func (e ErrVersionNotPublished) Error() string {
	return "version not published"
}
