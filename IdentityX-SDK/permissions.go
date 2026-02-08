package goauth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID         uuid.UUID        `json:"id"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	Object     string           `json:"object"`
	Action     string           `json:"action"`
	Conditions *json.RawMessage `json:"conditions"`
	CreatedAt  time.Time        `json:"created_at"`
}

type PermissionService struct {
	client *Client
}

type PermissionBuilder interface {
	Object(obj any) PermissionBuilder
	Action(act any) PermissionBuilder
	Conditions(cond any) PermissionBuilder
	Create(ctx context.Context) (*Permission, error)
}

type permissionBuilder struct {
	client     *Client
	object     string
	action     string
	conditions any
}

func (s *PermissionService) Define() PermissionBuilder {
	return &permissionBuilder{
		client: s.client,
	}
}

func (b *permissionBuilder) Object(obj any) PermissionBuilder {
	switch v := obj.(type) {
	case string:
		b.object = v
	case FinalizedBuilder:
		b.object = v.String()
	}
	return b
}

func (b *permissionBuilder) Action(act any) PermissionBuilder {
	switch v := act.(type) {
	case string:
		b.action = v
	case FinalizedBuilder:
		b.action = v.String()
	}
	return b
}

func (b *permissionBuilder) Conditions(cond any) PermissionBuilder {
	if cb, ok := cond.(ConditionBuilder); ok {
		b.conditions = cb.Build()
	} else {
		b.conditions = cond
	}
	return b
}

func (b *permissionBuilder) Create(ctx context.Context) (*Permission, error) {
	if err := validateObject(b.object); err != nil {
		return nil, err
	}
	if err := validateAction(b.action); err != nil {
		return nil, err
	}

	reqBody := map[string]any{
		"object":     b.object,
		"action":     b.action,
		"conditions": b.conditions,
	}
	path := fmt.Sprintf("/projects/%s/permissions", b.client.projectID)
	req, err := b.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data Permission `json:"data"`
	}
	err = b.client.do(req, &res)
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