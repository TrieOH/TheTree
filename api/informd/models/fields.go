package models

import (
	"encoding/json"
	"time"
	"uuid"
)

type FieldType string

const (
	FieldTypeString   FieldType = "string"
	FieldTypeEmail    FieldType = "email"
	FieldTypeInt      FieldType = "int"
	FieldTypeFloat    FieldType = "float"
	FieldTypeBool     FieldType = "bool"
	FieldTypeDate     FieldType = "date"
	FieldTypeTime     FieldType = "time"
	FieldTypeDatetime FieldType = "datetime"
	FieldTypeSelect   FieldType = "select"
	FieldTypeFile     FieldType = "file"
	FieldTypePhone    FieldType = "phone"
	FieldTypeURL      FieldType = "url"
)

type SelectBehaviour string

const (
	SelectBehaviourCheckbox         SelectBehaviour = "checkbox"
	SelectBehaviourRadio            SelectBehaviour = "radio"
	SelectBehaviourDropdownCheckbox SelectBehaviour = "dropdown-checkbox"
	SelectBehaviourDropdownRadio    SelectBehaviour = "dropdown-radio"
)

type SelectValueType string

const (
	SelectValueTypeString   SelectValueType = "string"
	SelectValueTypeEmail    SelectValueType = "email"
	SelectValueTypeInt      SelectValueType = "int"
	SelectValueTypeFloat    SelectValueType = "float"
	SelectValueTypeDate     SelectValueType = "date"
	SelectValueTypeTime     SelectValueType = "time"
	SelectValueTypeDatetime SelectValueType = "datetime"
	SelectValueTypePhone    SelectValueType = "phone"
	SelectValueTypeURL      SelectValueType = "url"
)

type Field struct {
	ID           uuid.UUID        `json:"id"`
	StepID       uuid.UUID        `json:"step_id"                 validate:"required"`
	Key          string           `json:"key"                     validate:"required"`
	Title        string           `json:"title"                   validate:"required"`
	Description  *string          `json:"description"`
	PositionHint int              `json:"position_hint"           validate:"required,gte=1"`
	Required     bool             `json:"required"`
	Type         FieldType        `json:"type"                    validate:"required"`
	Placeholder  *json.RawMessage `json:"placeholder,omitempty"`
	DefaultValue *json.RawMessage `json:"default_value,omitempty"`
	Config       *json.RawMessage `json:"config,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type FieldSelectConfig struct {
	FieldID   uuid.UUID       `json:"field_id"`
	Behaviour SelectBehaviour `json:"behaviour"  validate:"required"`
	ValueType SelectValueType `json:"value_type" validate:"required"`
	Options   json.RawMessage `json:"options"    validate:"required"`
}

// CreateFieldRequest is the HTTP request body for creating a field.
// SelectConfig is only required when Type is "select".

type CreateFieldSelectConfigInput struct {
	Behaviour SelectBehaviour `json:"behaviour"  validate:"required"`
	ValueType SelectValueType `json:"value_type" validate:"required"`
	Options   json.RawMessage `json:"options"    validate:"required"`
}

type CreateStepFieldInput struct {
	FormID       uuid.UUID
	StepID       uuid.UUID
	Key          string
	Title        string
	Description  *string
	PositionHint int
	Required     bool
	Type         FieldType
	Placeholder  *json.RawMessage
	DefaultValue *json.RawMessage
	Config       *json.RawMessage
	SelectConfig *CreateFieldSelectConfigInput
}

type UpdateStepFieldInput struct {
	StepID       uuid.UUID
	ID           uuid.UUID
	Key          string
	Title        string
	Description  *string
	PositionHint int
	Required     bool
	Type         FieldType
	Placeholder  *json.RawMessage
	DefaultValue *json.RawMessage
	Config       *json.RawMessage
	SelectConfig *CreateFieldSelectConfigInput
}

func UpdateStepFieldInputToField(i UpdateStepFieldInput) Field {
	return Field{
		ID:           i.ID,
		StepID:       i.StepID,
		Key:          i.Key,
		Title:        i.Title,
		Description:  i.Description,
		PositionHint: i.PositionHint,
		Required:     i.Required,
		Type:         i.Type,
		Placeholder:  i.Placeholder,
		DefaultValue: i.DefaultValue,
		Config:       i.Config,
	}
}
