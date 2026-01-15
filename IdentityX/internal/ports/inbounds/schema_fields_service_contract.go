package inbounds

import (
	"GoAuth/internal/domain/field"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SchemaFieldInput struct {
	SchemaID      uuid.UUID
	ProjectID     uuid.UUID
	VersionNumber int
	Fields        []InputField
}

type InputField struct {
	Key             string
	SchemaID        uuid.UUID
	SchemaVersionID uuid.UUID
	Type            string
	Owner           string
	Title           string
	Description     *string
	Placeholder     *string
	Required        bool
	Mutable         bool
	DefaultValue    *json.RawMessage
	Position        int
}

type OutputField struct {
	ObjectID        uuid.UUID
	ID              uuid.UUID
	Key             string
	SchemaID        uuid.UUID
	SchemaVersionID uuid.UUID
	Type            string
	Owner           string
	Title           string
	Description     *string
	Placeholder     *string
	Required        bool
	Mutable         bool
	DefaultValue    *json.RawMessage
	Position        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func FieldSliceToOutputFieldSlice(fs []field.Field) []OutputField {
	out := make([]OutputField, 0, len(fs))
	for _, f := range fs {
		out = append(out, *FieldToOutputField(&f))
	}
	return out
}

func FieldToOutputField(f *field.Field) *OutputField {
	if f == nil {
		return nil
	}
	return &OutputField{
		ObjectID:        f.ObjectID,
		ID:              f.ID,
		Key:             f.Key,
		SchemaID:        f.SchemaID,
		SchemaVersionID: f.SchemaVersionID,
		Type:            string(f.Type),
		Owner:           string(f.Owner),
		Title:           f.Title,
		Description:     f.Description,
		Placeholder:     f.Placeholder,
		Required:        f.Required,
		Mutable:         f.Mutable,
		DefaultValue:    f.DefaultValue,
		Position:        f.Position,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
	}
}

type ErrSchemaVersionMismatchLatest struct{}

func (e ErrSchemaVersionMismatchLatest) Error() string {
	return "version number does not match latest version"
}

type ErrAddFieldsToNonDraftVersion struct{}

func (e ErrAddFieldsToNonDraftVersion) Error() string {
	return "cannot add fields to a non-draft version"
}

type ErrInvalidFieldType struct {
	Type string
	Key  string
}

func (e ErrInvalidFieldType) Error() string {
	return "invalid field type (" + e.Type + ") for field: " + e.Key
}

type ErrInvalidFieldOwner struct {
	Owner string
	Key   string
}

func (e ErrInvalidFieldOwner) Error() string {
	return "invalid owner type (" + e.Owner + ") for field: " + e.Key
}
