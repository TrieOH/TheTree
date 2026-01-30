package dto

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

func SchemaOutputSliceToResponse(out []inbounds.SchemaOutput) []SchemaResponse {
	if out == nil {
		return nil
	}

	outSlice := make([]SchemaResponse, 0, len(out))
	for _, schema := range out {
		outSlice = append(outSlice, SchemaResponse{
			ID:               schema.ID,
			ProjectID:        schema.ProjectID,
			Title:            schema.Title,
			FlowID:           schema.FlowID,
			Type:             schema.Type,
			CurrentVersionID: schema.CurrentVersionID,
			Status:           schema.Status,
			CreatedAt:        schema.CreatedAt,
			UpdatedAt:        schema.UpdatedAt,
		})
	}
	return outSlice
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

	baseResponse := SchemaOutputToResponse(&out.SchemaOutput)
	schemaDTO := &VerboseSchemaResponse{
		SchemaResponse: *baseResponse,
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
				ID:               version.ID,
				SchemaID:         version.SchemaID,
				BasedOnVersionID: version.BasedOnVersionID,
				VersionNumber:    version.VersionNumber,
				Status:           string(version.Status),
				CreatedAt:        version.CreatedAt,
				UpdatedAt:        version.UpdatedAt,
			},
			Fields: fields,
		}
		versionsDTO = append(versionsDTO, versionOutput)
	}
	schemaDTO.Versions = versionsDTO

	return schemaDTO
}

// FormResponse represents a ready-to-render form with full field configuration
type FormResponse struct {
	ID            string      `json:"id"`
	SchemaID      string      `json:"schema_id"`
	Title         string      `json:"title"`
	FlowID        string      `json:"flow_id"`
	SchemaType    string      `json:"schema_type"`
	VersionID     string      `json:"version_id"`
	VersionNumber int         `json:"version_number"`
	Status        string      `json:"status"`
	Fields        []FormField `json:"fields"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type FormField struct {
	ID              string        `json:"id"`
	ObjectID        string        `json:"object_id"`
	Key             string        `json:"key"`
	Type            string        `json:"type"`
	Owner           string        `json:"owner"`
	Title           string        `json:"title"`
	Description     *string       `json:"description"`
	Placeholder     *string       `json:"placeholder"`
	Required        bool          `json:"required"`
	Mutable         bool          `json:"mutable"`
	DefaultValue    interface{}   `json:"default_value"`
	Position        int           `json:"position"`
	Options         []FieldOption `json:"options"`
	VisibilityRules []FieldRule   `json:"visibility_rules"`
	RequiredRules   []FieldRule   `json:"required_rules"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type FieldOption struct {
	ID       string `json:"id"`
	Value    string `json:"value"`
	Label    string `json:"label"`
	Position int    `json:"position"`
}

type FieldRule struct {
	ID               string      `json:"id"`
	DependsOnFieldID string      `json:"depends_on_field_id"`
	Operator         string      `json:"operator"`
	Value            interface{} `json:"value"`
}

func FormOutputToFormToResponse(form *inbounds.FormOutput) FormResponse {
	response := FormResponse{
		ID:            form.VersionID.String(),
		SchemaID:      form.SchemaID.String(),
		Title:         form.Title,
		FlowID:        form.FlowID,
		SchemaType:    form.SchemaType,
		VersionID:     form.VersionID.String(),
		VersionNumber: form.VersionNumber,
		Status:        form.Status,
		CreatedAt:     form.CreatedAt,
		UpdatedAt:     form.UpdatedAt,
		Fields:        make([]FormField, 0, len(form.Fields)),
	}

	for _, f := range form.Fields {
		field := FormField{
			ID:              f.ID.String(),
			ObjectID:        f.ObjectID.String(),
			Key:             f.Key,
			Type:            f.Type,
			Owner:           f.Owner,
			Title:           f.Title,
			Description:     f.Description,
			Placeholder:     f.Placeholder,
			Required:        f.Required,
			Mutable:         f.Mutable,
			Position:        f.Position,
			CreatedAt:       f.CreatedAt,
			UpdatedAt:       f.UpdatedAt,
			Options:         make([]FieldOption, 0),
			VisibilityRules: make([]FieldRule, 0),
			RequiredRules:   make([]FieldRule, 0),
		}

		if f.DefaultValue != nil {
			var val interface{}
			if err := json.Unmarshal(*f.DefaultValue, &val); err != nil {
				logs.L().Error("error unmarshalling default value", zap.Error(err))
			}
			field.DefaultValue = val
		}

		for _, opt := range f.Options {
			field.Options = append(field.Options, FieldOption{
				ID:       opt.ID.String(),
				Value:    opt.Value,
				Label:    opt.Label,
				Position: opt.Position,
			})
		}

		for _, rule := range f.VisibilityRules {
			r := FieldRule{
				ID:               rule.ID.String(),
				DependsOnFieldID: rule.DependsOnFieldID.String(),
				Operator:         rule.Operator,
			}
			if rule.Value != nil {
				var val interface{}
				if err := json.Unmarshal(*rule.Value, &val); err != nil {
					logs.L().Error("error unmarshalling rule value", zap.Error(err))
				}
				r.Value = val
			}
			field.VisibilityRules = append(field.VisibilityRules, r)
		}

		for _, rule := range f.RequiredRules {
			r := FieldRule{
				ID:               rule.ID.String(),
				DependsOnFieldID: rule.DependsOnFieldID.String(),
				Operator:         rule.Operator,
			}
			if rule.Value != nil {
				var val interface{}
				if err := json.Unmarshal(*rule.Value, &val); err != nil {
					logs.L().Error("error unmarshalling rule value", zap.Error(err))
				}
				r.Value = val
			}
			field.RequiredRules = append(field.RequiredRules, r)
		}

		response.Fields = append(response.Fields, field)
	}

	return response
}
