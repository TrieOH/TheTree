package goauth_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	baseURL := os.Getenv("GOAUTH_BASE_URL")
	apiKey := os.Getenv("GOAUTH_API_KEY")
	projectIDStr := os.Getenv("GOAUTH_PROJECT_ID")
	userIDStr := os.Getenv("GOAUTH_USER_ID")

	if baseURL == "" || apiKey == "" || projectIDStr == "" || userIDStr == "" {
		t.Skip("Skipping integration test: missing environment variables")
	}

	projectID, err := uuid.Parse(projectIDStr)
	require.NoError(t, err)

	userID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	ctx := context.Background()
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		ProjectID: projectID,
	})
	require.NoError(t, err)

	uniqueSuffix := uuid.New().String()[:8]
	objStr := "document" + uniqueSuffix
	roleName := "Tester-" + uniqueSuffix

	var permID uuid.UUID
	var roleID uuid.UUID
	var scopeNilExtID uuid.UUID
	var scopeWithExtID uuid.UUID

	t.Run("Setup Resources", func(t *testing.T) {
		// 1. Create Permission with simple string object/action
		perm, err := client.Permissions.Create(ctx, objStr, "read")
		require.NoError(t, err)
		permID = perm.ID

		// 2. Create Role
		role, err := client.Roles.Define(roleName).Create(ctx)
		require.NoError(t, err)
		roleID = role.ID

		// 3. Add Permission to Role
		err = client.Roles.AddPermission(ctx, roleID, permID)
		require.NoError(t, err)

		// 4. Create Scope with nil ExternalID
		s1, err := client.Scopes.Create(ctx, "ScopeNilExt-"+uniqueSuffix, nil)
		require.NoError(t, err)
		scopeNilExtID = s1.ID

		// 5. Create Scope with ExternalID
		extID := "ext-" + uniqueSuffix
		s2, err := client.Scopes.Create(ctx, "ScopeWithExt-"+uniqueSuffix, &extID)
		require.NoError(t, err)
		scopeWithExtID = s2.ID
	})

	t.Run("Global Role-Based Access", func(t *testing.T) {
		err = client.Roles.GiveToUser(ctx, userID, roleID, nil)
		require.NoError(t, err)
		defer client.Roles.TakeFromUser(ctx, userID, roleID, nil)

		eff, err := client.Permissions.GetEffective(ctx, userID, nil)
		require.NoError(t, err)
		assert.True(t, containsPermission(eff, permID))

		allowed, err := client.Authz.Check().User(userID).Object(objStr).Action("read").Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Scoped Role-Based Access (Nil External ID)", func(t *testing.T) {
		err = client.Roles.GiveToUser(ctx, userID, roleID, &scopeNilExtID)
		require.NoError(t, err)
		defer client.Roles.TakeFromUser(ctx, userID, roleID, &scopeNilExtID)

		eff, err := client.Permissions.GetEffective(ctx, userID, &scopeNilExtID)
		require.NoError(t, err)
		assert.True(t, containsPermission(eff, permID))

		allowed, err := client.Authz.Check().User(userID).Scope(scopeNilExtID).Object(objStr).Action("read").Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Scoped Role-Based Access (With External ID)", func(t *testing.T) {
		err = client.Roles.GiveToUser(ctx, userID, roleID, &scopeWithExtID)
		require.NoError(t, err)
		defer client.Roles.TakeFromUser(ctx, userID, roleID, &scopeWithExtID)

		eff, err := client.Permissions.GetEffective(ctx, userID, &scopeWithExtID)
		require.NoError(t, err)
		assert.True(t, containsPermission(eff, permID))

		allowed, err := client.Authz.Check().User(userID).Scope(scopeWithExtID).Object(objStr).Action("read").Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Direct Permission Access (Global)", func(t *testing.T) {
		err = client.Permissions.GiveDirect(ctx, userID, permID, nil)
		require.NoError(t, err)
		defer client.Permissions.TakeDirect(ctx, userID, permID, nil)

		eff, err := client.Permissions.GetEffective(ctx, userID, nil)
		require.NoError(t, err)
		assert.True(t, containsPermission(eff, permID))

		allowed, err := client.Authz.Check().User(userID).Object(objStr).Action("read").Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Direct Permission Access (Scoped)", func(t *testing.T) {
		err = client.Permissions.GiveDirect(ctx, userID, permID, &scopeNilExtID)
		require.NoError(t, err)
		defer client.Permissions.TakeDirect(ctx, userID, permID, &scopeNilExtID)

		eff, err := client.Permissions.GetEffective(ctx, userID, &scopeNilExtID)
		require.NoError(t, err)
		assert.True(t, containsPermission(eff, permID))

		allowed, err := client.Authz.Check().User(userID).Scope(scopeNilExtID).Object(objStr).Action("read").Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Exact Match Check", func(t *testing.T) {
		// Create permission with exact object/action
		exactPerm, err := client.Permissions.Create(ctx, "report", "generate")
		require.NoError(t, err)

		err = client.Permissions.GiveDirect(ctx, userID, exactPerm.ID, nil)
		require.NoError(t, err)
		defer client.Permissions.TakeDirect(ctx, userID, exactPerm.ID, nil)

		// Exact match should work
		allowed, err := client.Authz.Check().User(userID).Object("report").Action("generate").Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed, "Exact match should be allowed")

		// Different action should not match
		allowed, err = client.Authz.Check().User(userID).Object("report").Action("delete").Allowed(ctx)
		require.NoError(t, err)
		assert.False(t, allowed, "Different action should not match")

		// Different object should not match
		allowed, err = client.Authz.Check().User(userID).Object("document").Action("generate").Allowed(ctx)
		require.NoError(t, err)
		assert.False(t, allowed, "Different object should not match")
	})

	t.Run("Permission Validation", func(t *testing.T) {
		// Invalid object with wildcard should fail
		_, err := client.Permissions.Create(ctx, "doc:*", "read")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid object format")

		// Invalid action with separator should fail
		_, err = client.Permissions.Create(ctx, "document", "read:write")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid action format")

		// Valid simple names should work
		validPerm, err := client.Permissions.Create(ctx, "user_profile", "update")
		require.NoError(t, err)
		assert.Equal(t, "user_profile", validPerm.Object)
		assert.Equal(t, "update", validPerm.Action)

		// Cleanup
		client.Permissions.TakeDirect(ctx, userID, validPerm.ID, nil)
	})

	t.Run("Project Discovery & User Management", func(t *testing.T) {
		// 1. Get JWKS
		jwks, err := client.Tokens.GetJWKS(ctx, false)
		require.NoError(t, err)
		assert.NotEmpty(t, jwks.Keys, "JWKS should contain at least one key")

		// 2. List Users
		users, err := client.Users.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, users, "Project should have at least one user (the provided GOAUTH_USER_ID)")

		found := false
		for _, u := range users {
			if u.ID == userID {
				found = true
				break
			}
		}
		assert.True(t, found, "The provided user ID should be in the project users list")

		// 3. Get User
		user, err := client.Users.Get(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, projectID, user.ProjectID)
	})

	t.Run("Resolution Cross-Check (Isolation)", func(t *testing.T) {
		err = client.Roles.GiveToUser(ctx, userID, roleID, &scopeNilExtID)
		require.NoError(t, err)
		defer client.Roles.TakeFromUser(ctx, userID, roleID, &scopeNilExtID)

		effA, _ := client.Permissions.GetEffective(ctx, userID, &scopeNilExtID)
		assert.True(t, containsPermission(effA, permID))

		effB, _ := client.Permissions.GetEffective(ctx, userID, &scopeWithExtID)
		assert.False(t, containsPermission(effB, permID), "Scope B should not see permissions from Scope A")

		effGlobal, _ := client.Permissions.GetEffective(ctx, userID, nil)
		assert.False(t, containsPermission(effGlobal, permID), "Global level should not see permissions from Scope A")
	})

	t.Run("List Permissions with Filters", func(t *testing.T) {
		// Create a few permissions
		perm1, _ := client.Permissions.Create(ctx, "file", "read")
		perm2, _ := client.Permissions.Create(ctx, "file", "write")
		perm3, _ := client.Permissions.Create(ctx, "database", "read")

		// List all
		allPerms, err := client.Permissions.List(ctx, "", "")
		require.NoError(t, err)
		assert.True(t, len(allPerms) >= 3)

		// List by object
		filePerms, err := client.Permissions.List(ctx, "file", "")
		require.NoError(t, err)
		assert.True(t, len(filePerms) >= 2)
		for _, p := range filePerms {
			assert.Equal(t, "file", p.Object)
		}

		// List by action
		readPerms, err := client.Permissions.List(ctx, "", "read")
		require.NoError(t, err)
		assert.True(t, len(readPerms) >= 2)
		for _, p := range readPerms {
			assert.Equal(t, "read", p.Action)
		}

		// List by both
		fileReadPerms, err := client.Permissions.List(ctx, "file", "read")
		require.NoError(t, err)
		assert.True(t, len(fileReadPerms) >= 1)
		for _, p := range fileReadPerms {
			assert.Equal(t, "file", p.Object)
			assert.Equal(t, "read", p.Action)
		}

		// Cleanup
		client.Permissions.TakeDirect(ctx, userID, perm1.ID, nil)
		client.Permissions.TakeDirect(ctx, userID, perm2.ID, nil)
		client.Permissions.TakeDirect(ctx, userID, perm3.ID, nil)
	})
}

func containsPermission(perms []goauth.Permission, id uuid.UUID) bool {
	for _, p := range perms {
		if p.ID == id {
			return true
		}
	}
	return false
}

func handleError(op string, err error) {
	if fe, ok := fail.As(err); ok {
		fmt.Printf("Error during %s: [%s] %s\n", op, fe.ID, fe.Message)
		if apiID, ok := fe.Meta["api_id"]; ok && apiID != "" {
			fmt.Printf("  API Error ID: %v\n", apiID)
		}
		if traces, ok := fe.Meta["api_traces"]; ok {
			if tSlice, ok := traces.([]string); ok && len(tSlice) > 0 {
				fmt.Println("  API Traces:")
				for _, t := range tSlice {
					fmt.Printf("    - %s\n", t)
				}
			}
		}
	} else {
		fmt.Printf("Error during %s: %v\n", op, err)
	}
}
