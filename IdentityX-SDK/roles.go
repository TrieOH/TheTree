package goauth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   *uuid.UUID `json:"project_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ScopeID     *uuid.UUID `json:"scope_id"`
	ScopeName   *string    `json:"scope_name"`
	ExternalID  *string    `json:"external_id"`
}

type RoleService struct {
	client *Client
}

type RoleBuilder interface {
	Description(desc string) RoleBuilder
	Create(ctx context.Context) (*Role, error)
}

type roleBuilder struct {
	client      *Client
	name        string
	description *string
}

func (s *RoleService) Define(name string) RoleBuilder {
	return &roleBuilder{
		client: s.client,
		name:   name,
	}
}

func (b *roleBuilder) Description(desc string) RoleBuilder {
	b.description = &desc
	return b
}

func (b *roleBuilder) Create(ctx context.Context) (*Role, error) {
	reqBody := map[string]any{
		"name":        b.name,
		"description": b.description,
	}
	path := fmt.Sprintf("/projects/%s/roles", b.client.projectID)
	req, err := b.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data Role `json:"data"`
	}
	err = b.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (s *RoleService) List(ctx context.Context) ([]Role, error) {
	path := fmt.Sprintf("/projects/%s/roles", s.client.projectID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data []Role `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (s *RoleService) Get(ctx context.Context, roleID uuid.UUID) (*Role, error) {
	path := fmt.Sprintf("/projects/%s/roles/%s", s.client.projectID, roleID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data Role `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (s *RoleService) AddPermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	path := fmt.Sprintf("/projects/%s/roles/%s/permissions/%s", s.client.projectID, roleID, permissionID)
	req, err := s.client.newRequest(ctx, "POST", path, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *RoleService) RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	path := fmt.Sprintf("/projects/%s/roles/%s/permissions/%s", s.client.projectID, roleID, permissionID)
	req, err := s.client.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *RoleService) GiveToUser(ctx context.Context, entityID, roleID uuid.UUID, scopeID *uuid.UUID) error {
	reqBody := map[string]any{
		"role_id":  roleID,
		"scope_id": scopeID,
	}
	path := fmt.Sprintf("/projects/%s/identities/%s/roles", s.client.projectID, entityID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *RoleService) TakeFromUser(ctx context.Context, entityID, roleID uuid.UUID, scopeID *uuid.UUID) error {
	reqBody := map[string]any{
		"role_id":  roleID,
		"scope_id": scopeID,
	}
	path := fmt.Sprintf("/projects/%s/identities/%s/roles", s.client.projectID, entityID)
	req, err := s.client.newRequest(ctx, "DELETE", path, reqBody)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *RoleService) Give(ctx context.Context, entityID uuid.UUID, roleName string, scopeID *uuid.UUID) error {
	reqBody := map[string]any{
		"role_name": roleName,
		"scope_id":  scopeID,
	}
	path := fmt.Sprintf("/projects/%s/identities/%s/roles/by-name", s.client.projectID, entityID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *RoleService) Take(ctx context.Context, entityID uuid.UUID, roleName string, scopeID *uuid.UUID) error {
	reqBody := map[string]any{
		"role_name": roleName,
		"scope_id":  scopeID,
	}
	path := fmt.Sprintf("/projects/%s/identities/%s/roles/by-name", s.client.projectID, entityID)
	req, err := s.client.newRequest(ctx, "DELETE", path, reqBody)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

type RoleDefinition struct {
	Name        string
	Permissions []PermissionDefinition
	Meta        map[string]interface{}
}

type EnsureRoleResult struct {
	Name    string `json:"name"`
	Created bool   `json:"created"`
}

func (s *RoleService) EnsureExists(ctx context.Context, roles []RoleDefinition) ([]EnsureRoleResult, error) {
	reqBody := map[string]any{
		"roles": roles,
	}
	path := fmt.Sprintf("/projects/%s/roles/ensure", s.client.projectID)
	req, err := s.client.newRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data struct {
			Roles []EnsureRoleResult `json:"roles"`
		} `json:"data"`
	}
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data.Roles, nil
}
