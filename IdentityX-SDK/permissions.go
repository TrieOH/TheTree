package goauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID *uuid.UUID `json:"project_id"`
	Object    string     `json:"object"`
	Action    string     `json:"action"`
	CreatedAt time.Time  `json:"created_at"`
}

type PermissionService struct {
	client *Client
}

func (s *PermissionService) Create(ctx context.Context, object, action string) (*Permission, error) {
	if err := validateObject(object); err != nil {
		return nil, err
	}
	if err := validateAction(action); err != nil {
		return nil, err
	}

	reqBody := map[string]any{
		"object": object,
		"action": action,
	}
	path := fmt.Sprintf("/projects/%s/permissions", s.client.projectID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data Permission `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (s *PermissionService) List(ctx context.Context, object, action string) ([]Permission, error) {
	path := fmt.Sprintf("/projects/%s/permissions", s.client.projectID)

	queryParams := make([]string, 0)
	if object != "" {
		if err := validateObject(object); err != nil {
			return nil, err
		}
		queryParams = append(queryParams, fmt.Sprintf("object=%s", object))
	}
	if action != "" {
		if err := validateAction(action); err != nil {
			return nil, err
		}
		queryParams = append(queryParams, fmt.Sprintf("action=%s", action))
	}

	if len(queryParams) > 0 {
		path = fmt.Sprintf("%s?%s", path, strings.Join(queryParams, "&"))
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data []Permission `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (s *PermissionService) GiveDirect(ctx context.Context, entityID uuid.UUID, permissionID uuid.UUID, scopeID *uuid.UUID) error {
	reqBody := map[string]any{
		"permission_id": permissionID,
		"scope_id":      scopeID,
	}
	path := fmt.Sprintf("/projects/%s/identities/%s/permissions", s.client.projectID, entityID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *PermissionService) TakeDirect(ctx context.Context, entityID uuid.UUID, permissionID uuid.UUID, scopeID *uuid.UUID) error {
	reqBody := map[string]any{
		"permission_id": permissionID,
		"scope_id":      scopeID,
	}
	path := fmt.Sprintf("/projects/%s/identities/%s/permissions", s.client.projectID, entityID)
	req, err := s.client.newRequest(ctx, "DELETE", path, reqBody)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *PermissionService) GetEffective(ctx context.Context, entityID uuid.UUID, scopeID *uuid.UUID) ([]Permission, error) {
	path := fmt.Sprintf("/projects/%s/identities/%s/permissions", s.client.projectID, entityID)
	if scopeID != nil {
		path = fmt.Sprintf("%s?scope_id=%s", path, scopeID)
	}
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data []Permission `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

type PermissionDefinition struct {
	Object string
	Action string
	Meta   map[string]interface{}
}

type EnsurePermissionResult struct {
	Object  string `json:"object"`
	Action  string `json:"action"`
	Created bool   `json:"created"`
}

func (s *PermissionService) EnsureExists(ctx context.Context, permissions []PermissionDefinition) ([]EnsurePermissionResult, error) {
	reqBody := map[string]any{
		"permissions": permissions,
	}
	path := fmt.Sprintf("/projects/%s/permissions/ensure", s.client.projectID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data struct {
			Permissions []EnsurePermissionResult `json:"permissions"`
		} `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data.Permissions, nil
}
