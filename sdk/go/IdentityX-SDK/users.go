package idx

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeClient  UserType = "client"
	UserTypeProject UserType = "project"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	UserType    UserType   `json:"user_type"`
	ProjectID   *uuid.UUID `json:"project_id"`
	Email       string     `json:"email"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	IsVerified  bool       `json:"is_verified"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
}

type UserService struct {
	client *Client
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
	var out []User
	err := s.client.DoRequest(ctx, "GET", fmt.Sprintf("/projects/%s/users", s.client.projectID), nil, &out)
	return out, err
}

func (s *UserService) Get(ctx context.Context, userID uuid.UUID) (*User, error) {
	var out User
	err := s.client.DoRequest(ctx, "GET", fmt.Sprintf("/projects/%s/users/%s", s.client.projectID, userID), nil, &out)
	return &out, err
}
