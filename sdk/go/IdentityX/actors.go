package idx

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthMethod string

const (
	PasswordAuthMethod AuthMethod = "password"
	APIKeyAuthMethod   AuthMethod = "api_key"
	GoogleAuthMethod   AuthMethod = "google_auth"
	GithubAuthMethod   AuthMethod = "github_auth"
)

type ActorType string

const (
	HumanActorType   ActorType = "human"
	ServiceActorType ActorType = "service"
	MachineActorType ActorType = "machine"
)

type Actor struct {
	ID         uuid.UUID        `json:"id"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	AuthMethod AuthMethod       `json:"auth_method"`
	VerifiedAt *time.Time       `json:"verified_at"`
	Email      *string          `json:"email"`
	Type       ActorType        `json:"type"`
	Metadata   *json.RawMessage `json:"metadata"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  *time.Time       `json:"deleted_at"`
}

type ActorService struct {
	client *Client
}

func (s *ActorService) GetByEmail(ctx context.Context, email string) (*Actor, error) {
	var res Actor
	req := s.client.http.R().SetContext(ctx).
		SetPathParam("project_id", s.client.projectID.String()).
		SetPathParam("email", email).
		SetResult(&res)
	err := s.client.send(req, http.MethodGet, "/projects/{project_id}/actors/{email}:by_email")
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ActorService) GetByID(ctx context.Context, id uuid.UUID) (*Actor, error) {
	var res Actor
	req := s.client.http.R().SetContext(ctx).
		SetPathParam("project_id", s.client.projectID.String()).
		SetPathParam("id", id.String()).
		SetResult(&res)
	err := s.client.send(req, http.MethodGet, "/projects/{project_id}/actors/{id}")
	if err != nil {
		return nil, err
	}
	return &res, nil
}
