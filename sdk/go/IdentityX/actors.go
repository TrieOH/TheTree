package idx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthMethod string

const (
	PasswordAuthMethod AuthMethod = "password"
	ApiKeyAuthMethod   AuthMethod = "api_key"
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
	path := fmt.Sprintf("/projects/%s/actors/%s:by_email", s.client.projectID, email)
	if err := s.client.DoRequest(ctx, "GET", path, nil, &res); err != nil {
		return nil, err // *sdkkit.SDKError or *sdkkit.APIError — both appropriate here
	}
	return &res, nil
}

func (s *ActorService) GetByID(ctx context.Context, id uuid.UUID) (*Actor, error) {
	var res Actor
	path := fmt.Sprintf("/projects/%s/actors/%s", s.client.projectID, id.String())
	if err := s.client.DoRequest(ctx, "GET", path, nil, &res); err != nil {
		return nil, err // *sdkkit.SDKError or *sdkkit.APIError — both appropriate here
	}
	return &res, nil
}
