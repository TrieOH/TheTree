package dto

import (
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/google/uuid"
)

type SchemaResponse struct {
	ID               uuid.UUID  `json:"id"`
	ProjectID        uuid.UUID  `json:"project_id"`
	Title            string     `json:"title"`
	FlowID           string     `json:"flow_id"`
	Type             string     `json:"type"`
	CurrentVersionID *uuid.UUID `json:"current_version_id"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func SchemaOutputToResponse(out *inbounds.SchemaOutput) *SchemaResponse {
	if out == nil {
		return nil
	}
	return &SchemaResponse{
		ID:               out.ID,
		ProjectID:        out.ProjectID,
		Title:            out.Title,
		FlowID:           out.FlowID,
		Type:             out.Type,
		CurrentVersionID: out.CurrentVersionID,
		Status:           out.Status,
		CreatedAt:        out.CreatedAt,
		UpdatedAt:        out.UpdatedAt,
	}
}

type VerboseSchemaResponse struct {
	SchemaResponse
	Versions []VersionVerboseResponse `json:"versions"`
}

func VerboseSchemaOutputToResponse(out *inbounds.SchemaVerboseOutput) *VerboseSchemaResponse {
	if out == nil {
		return nil
	}

	schemaDTO := &VerboseSchemaResponse{
		SchemaResponse: SchemaResponse{
			ID:               out.ID,
			ProjectID:        out.ProjectID,
			Title:            out.Title,
			FlowID:           out.FlowID,
			Type:             out.Type,
			CurrentVersionID: out.CurrentVersionID,
			Status:           out.Status,
			CreatedAt:        out.CreatedAt,
			UpdatedAt:        out.UpdatedAt,
		},
	}

	versionsDTO := make([]VersionVerboseResponse, 0, len(out.Versions))
	for _, version := range out.Versions {
		fields := make([]FieldResponse, 0, len(version.Fields))
		for _, f := range version.Fields {
			fields = append(fields, FieldResponse{
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
		versionOutput := VersionVerboseResponse{
			SchemaVersionResponse: SchemaVersionResponse{
				ID:            version.ID,
				SchemaID:      version.SchemaID,
				VersionNumber: version.VersionNumber,
				Status:        string(version.Status),
				CreatedAt:     version.CreatedAt,
				UpdatedAt:     version.UpdatedAt,
			},
			Fields: fields,
		}
		versionsDTO = append(versionsDTO, versionOutput)
	}
	schemaDTO.Versions = versionsDTO

	return schemaDTO
}
