package idx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProjectUser struct {
	ID          uuid.UUID       `json:"id"`
	ProjectID   uuid.UUID       `json:"project_id"`
	Email       string          `json:"email"`
	UserType    string          `json:"user_type"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastLoginAt *time.Time      `json:"last_login_at,omitempty"`
	IsVerified  bool            `json:"is_verified"`
	VerifiedAt  *time.Time      `json:"verified_at,omitempty"`
}

type UserService struct {
	client *Client
}

func (s *UserService) List(ctx context.Context) ([]ProjectUser, error) {
	path := fmt.Sprintf("/projects/%s/users", s.client.projectID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data []ProjectUser `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (s *UserService) Get(ctx context.Context, userID uuid.UUID) (*ProjectUser, error) {
	path := fmt.Sprintf("/projects/%s/users/%s", s.client.projectID, userID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data ProjectUser `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}
