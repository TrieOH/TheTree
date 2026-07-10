package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProjectProfileSchema holds the JSON Schema used to validate actor profiles.
// When project_id is nil, the schema acts as the platform-wide default.
type ProjectProfileSchema struct {
	ProjectID *uuid.UUID      `json:"project_id"`
	Schema    json.RawMessage `json:"schema"`
	Version   int             `json:"version"`
	Active    bool            `json:"active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// UpsertProfileSchemaRequest is the HTTP request body for setting a profile schema.
type UpsertProfileSchemaRequest struct {
	Schema json.RawMessage `json:"schema"`
	Active bool            `json:"active"`
}

func (r UpsertProfileSchemaRequest) ToInput(projectID *uuid.UUID) UpsertProfileSchemaInput {
	return UpsertProfileSchemaInput{
		ProjectID: projectID,
		Schema:    r.Schema,
		Active:    r.Active,
	}
}

type UpsertProfileSchemaInput struct {
	ProjectID *uuid.UUID
	Schema    json.RawMessage
	Active    bool
}

// UpsertProfileRequest is the HTTP request body for upserting an actor's profile.
type UpsertProfileRequest struct {
	Profile json.RawMessage `json:"profile"`
}

func (r UpsertProfileRequest) ToInput(actorID uuid.UUID) UpsertProfileInput {
	return UpsertProfileInput{
		ActorID: actorID,
		Profile: r.Profile,
	}
}

type UpsertProfileInput struct {
	ActorID uuid.UUID
	Profile json.RawMessage
}
