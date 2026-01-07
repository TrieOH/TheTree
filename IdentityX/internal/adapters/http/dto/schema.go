package dto

import (
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/google/uuid"
)

type DraftSchemaRequest struct {
	SchemaType string `json:"schema_type" validate:"required,oneof=core context sub-context"`
	Title      string `json:"title" validate:"required,max=255"`
	FlowID     string `json:"flow_id" validate:"required,max=63"`
}

type DraftSchemaResponse struct {
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

func SchemaOutputToResponse(out *inbounds.SchemaOutput) *DraftSchemaResponse {
	if out == nil {
		return nil
	}
	return &DraftSchemaResponse{
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

type VerboseSchemaResponse struct {
	DraftSchemaResponse
	Versions []VersionVerboseResponse `json:"versions"`
}

func VerboseSchemaOutputToResponse(out *inbounds.SchemaVerboseOutput) *VerboseSchemaResponse {
	if out == nil {
		return nil
	}

	schemaDTO := &VerboseSchemaResponse{
		DraftSchemaResponse: DraftSchemaResponse{
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
		versionOutput := VersionVerboseResponse{
			DraftSchemaVersionResponse: DraftSchemaVersionResponse{
				ID:            version.ID,
				SchemaID:      version.SchemaID,
				VersionNumber: version.VersionNumber,
				Status:        string(version.Status),
				CreatedAt:     version.CreatedAt,
				UpdatedAt:     version.UpdatedAt,
			},
			Fields: nil,
		}
		versionsDTO = append(versionsDTO, versionOutput)
	}

	schemaDTO.Versions = versionsDTO

	for i := range schemaDTO.Versions {
		for _, f := range out.Versions[i].Fields {
			schemaDTO.Versions[i].Fields = append(schemaDTO.Versions[i].Fields, FieldResponse{
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
	}

	return schemaDTO
}
