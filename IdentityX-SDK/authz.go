package goauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type AuthzService struct {
	client *Client
}

type CheckBuilder interface {
	User(id uuid.UUID) CheckBuilder
	Scope(id uuid.UUID) CheckBuilder
	Object(obj string) CheckBuilder
	Action(act string) CheckBuilder
	Allowed(ctx context.Context) (bool, error)
}

type checkBuilder struct {
	client  *Client
	userID  *uuid.UUID
	scopeID *uuid.UUID
	object  string
	action  string
}

func (s *AuthzService) Check() CheckBuilder {
	return &checkBuilder{client: s.client}
}

func (b *checkBuilder) User(id uuid.UUID) CheckBuilder {
	b.userID = &id
	return b
}

func (b *checkBuilder) Scope(id uuid.UUID) CheckBuilder {
	b.scopeID = &id
	return b
}

func (b *checkBuilder) Object(obj string) CheckBuilder {
	b.object = obj
	return b
}

func (b *checkBuilder) Action(act string) CheckBuilder {
	b.action = act
	return b
}

func (b *checkBuilder) Allowed(ctx context.Context) (bool, error) {
	if b.userID == nil {
		return false, fail.New(SDKMissingUserID).Trace("user ID required for authz check")
	}

	if err := validateObject(b.object); err != nil {
		return false, err
	}
	if err := validateAction(b.action); err != nil {
		return false, err
	}

	reqBody := map[string]any{
		"project_id": b.client.projectID,
		"scope_id":   b.scopeID,
		"entity_id":  b.userID,
		"object":     b.object,
		"action":     b.action,
	}

	req, err := b.client.newRequest(ctx, "POST", "/authz/check", reqBody)
	if err != nil {
		return false, err
	}

	resp, err := b.client.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		return false, b.client.handleErrorResponse(resp, body)
	}

	var res struct {
		Data struct {
			Allowed bool `json:"allowed"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return false, err
	}

	return res.Data.Allowed, nil
}
