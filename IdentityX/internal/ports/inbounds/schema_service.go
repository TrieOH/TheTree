package inbounds

import (
	"GoAuth/internal/domain/schema"
	"context"
	"time"

	"github.com/google/uuid"
)

type SchemaService interface {
	// Draft /projects/{id}/schemas
	Draft(ctx context.Context, in DraftSchemaInput) (*DraftSchemaOutput, error)
	// Publish /projects/{id}/schemas/{schemaID}/publish?flowID=xxx
	Publish(ctx context.Context, in PublishSchemaInput) error
}

type DraftSchemaInput struct {
	SchemaType string
	Title      string
	FlowID     string
	ProjectID  string
}

type PublishSchemaInput struct {
	FlowID    string
	ProjectID string
}

type DraftSchemaOutput struct {
	ID               uuid.UUID
	ProjectID        uuid.UUID
	Title            string
	FlowID           string
	Type             schema.Type
	CurrentVersionID *uuid.UUID
	Status           schema.Status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
