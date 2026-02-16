package goauth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Scope struct {
	ID         uuid.UUID  `json:"id"`
	ParentID   *uuid.UUID `json:"parent_id"`
	ProjectID  *uuid.UUID `json:"project_id"`
	ExternalID *string    `json:"external_id"`
	Type       string     `json:"type"`
	Name       *string    `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ScopeService struct {
	client *Client
}

func (s *ScopeService) Create(ctx context.Context, name string, externalID *string) (*Scope, error) {
	return s.CreateWithParent(ctx, name, externalID, nil)
}

func (s *ScopeService) CreateWithParent(ctx context.Context, name string, externalID *string, parentID *uuid.UUID) (*Scope, error) {
	reqBody := map[string]any{
		"name":        name,
		"external_id": externalID,
		"parent_id":   parentID,
	}
	path := fmt.Sprintf("/projects/%s/scopes", s.client.projectID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data Scope `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (s *ScopeService) List(ctx context.Context) ([]Scope, error) {
	path := fmt.Sprintf("/projects/%s/scopes", s.client.projectID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data []Scope `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (s *ScopeService) Get(ctx context.Context, scopeID uuid.UUID) (*Scope, error) {
	path := fmt.Sprintf("/projects/%s/scopes/%s", s.client.projectID, scopeID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data Scope `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}
