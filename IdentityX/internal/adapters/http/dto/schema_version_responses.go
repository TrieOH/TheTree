package dto

import (
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/google/uuid"
)

type SchemaVersionResponse struct {
	ID               uuid.UUID  `json:"id"`
	SchemaID         uuid.UUID  `json:"schema_id"`
	BasedOnVersionID *uuid.UUID `json:"based_on_version_id"`
	VersionNumber    int        `json:"version_number"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type VersionVerboseResponse struct {
	SchemaVersionResponse
	Fields []FieldResponse `json:"fields"`
}

func SchemaVersionOutputToResponse(out *inbounds.SchemaVersionOutput) *SchemaVersionResponse {
	if out == nil {
		return nil
	}
	return &SchemaVersionResponse{
		ID:               out.ID,
		SchemaID:         out.SchemaID,
		BasedOnVersionID: out.BasedOnVersionID,
		VersionNumber:    out.VersionNumber,
		Status:           string(out.Status),
		CreatedAt:        out.CreatedAt,
		UpdatedAt:        out.UpdatedAt,
	}
}

func VerboseVersionOutputToResponse(out *inbounds.VersionVerboseOutput) *VersionVerboseResponse {
	if out == nil {
		return nil
	}

	res := &VersionVerboseResponse{
		SchemaVersionResponse: SchemaVersionResponse{
			ID:               out.ID,
			SchemaID:         out.SchemaID,
			BasedOnVersionID: out.BasedOnVersionID,
			VersionNumber:    out.VersionNumber,
			Status:           string(out.Status),
			CreatedAt:        out.CreatedAt,
			UpdatedAt:        out.UpdatedAt,
		},
		Fields: nil,
	}

	res.Fields = make([]FieldResponse, 0, len(out.Fields))
	for _, f := range out.Fields {
		res.Fields = append(res.Fields, FieldResponse{
			ObjectID:        f.ObjectID,
			ID:              f.ID,
			Key:             f.Key,
			SchemaID:        f.SchemaID,
			SchemaVersionID: f.SchemaVersionID,
			Type:            f.Type,
			Owner:           f.Owner,
			Title:           f.Title,
			Description:     f.Description,
			Placeholder:     f.Placeholder,
			Required:        f.Required,
			Mutable:         f.Mutable,
			DefaultValue:    f.DefaultValue,
			Position:        f.Position,
			CreatedAt:       f.CreatedAt,
			UpdatedAt:       f.UpdatedAt,
		})
	}
	return res
}
