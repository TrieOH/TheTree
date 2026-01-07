package field

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	String   Type = "string"
	Int      Type = "int"
	Select   Type = "select"
	Radio    Type = "radio"
	Checkbox Type = "checkbox"
	Bool     Type = "bool"
)

type Owner string

const (
	System Owner = "system"
	Admin  Owner = "admin"
	User   Owner = "user"
)

type Field struct {
	ObjectID        uuid.UUID
	ID              uuid.UUID
	SchemaID        uuid.UUID
	SchemaVersionID uuid.UUID
	Key             string
	Type            Type
	Owner           Owner
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

type Option struct {
	ID       uuid.UUID
	FieldID  uuid.UUID
	Value    string
	Label    string
	Position int
}

type RuleOperator string

const (
	RuleOperatorEquals    RuleOperator = "equals"
	RuleOperatorNotEquals RuleOperator = "not_equals"
	RuleOperatorIn        RuleOperator = "in"
	RuleOperatorNotIn     RuleOperator = "not_in"
	RuleOperatorExists    RuleOperator = "exists"
	RuleOperatorNotExists RuleOperator = "not_exists"
)

type RequiredRule struct {
	ID               uuid.UUID
	FieldID          uuid.UUID
	DependsOnFieldID uuid.UUID
	Operator         RuleOperator
	Value            *json.RawMessage
	CreatedAt        time.Time
}

type VisibilityRule struct {
	ID               uuid.UUID
	FieldID          uuid.UUID
	DependsOnFieldID uuid.UUID
	Operator         RuleOperator
	Value            *json.RawMessage
	CreatedAt        time.Time
}
