package schema

import (
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Type string

const (
	Core       Type = "core"
	Context    Type = "context"
	SubContext Type = "sub-context"
)

func IsValidSchemaType(s string) bool {
	switch Type(s) {
	case Core, Context, SubContext:
		return true
	default:
		return false
	}
}

type ReservedFlowID string

const NoFlowID ReservedFlowID = "none"

func IsFlowIDReserved(flowID string) bool {
	switch ReservedFlowID(flowID) {
	case NoFlowID:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Schema struct {
	ID               uuid.UUID
	ProjectID        uuid.UUID
	Title            string
	FlowID           string
	Type             Type
	CurrentVersionID *uuid.UUID
	Status           Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (s Schema) IsVersion(versionID uuid.UUID) bool {
	return s.CurrentVersionID != nil && *s.CurrentVersionID == versionID
}

type DiffResult struct {
	FieldsChanged          *bool `json:"fields_changed"`
	OptionsChanged         *bool `json:"options_changed"`
	VisibilityRulesChanged *bool `json:"visibility_rules_changed"`
	RequiredRulesChanged   *bool `json:"required_rules_changed"`
}

func (r DiffResult) HasAnyChanges() bool {
	return isTrue(r.FieldsChanged) ||
		isTrue(r.OptionsChanged) ||
		isTrue(r.VisibilityRulesChanged) ||
		isTrue(r.RequiredRulesChanged)
}

func isTrue(b *bool) bool {
	return b != nil && *b
}

func (r DiffResult) Annotate(span trace.Span) {
	if r.FieldsChanged != nil {
		span.SetAttributes(attribute.Bool("fields_changed", *r.FieldsChanged))
	}
	if r.OptionsChanged != nil {
		span.SetAttributes(attribute.Bool("options_changed", *r.OptionsChanged))
	}
	if r.VisibilityRulesChanged != nil {
		span.SetAttributes(attribute.Bool("visibility_rules_changed", *r.VisibilityRulesChanged))
	}
	if r.RequiredRulesChanged != nil {
		span.SetAttributes(attribute.Bool("required_rules_changed", *r.RequiredRulesChanged))
	}
	span.SetAttributes(attribute.Bool("has_any_changes", r.HasAnyChanges()))
}
