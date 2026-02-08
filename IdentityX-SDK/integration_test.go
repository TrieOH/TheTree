package goauth_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"github.com/MintzyG/fail/v3"
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
	objStr := "doc:test" + uniqueSuffix
	roleName := "Tester-" + uniqueSuffix

	var permID uuid.UUID
	var roleID uuid.UUID
	var scopeNilExtID uuid.UUID
	var scopeWithExtID uuid.UUID

	t.Run("Setup Resources", func(t *testing.T) {
		// 1. Create Permission
		perm, err := client.Permissions.Define().
			Object(objStr).
			Action(goauth.Action("read")).
			Create(ctx)
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

	t.Run("Complex Patterns Check", func(t *testing.T) {
		complexObj, _ := goauth.Object("events", "test"+uniqueSuffix).NamespaceAny("users").Build()
		complexAct, _ := goauth.Action("attendance").Any().Sub("edit").Build()
		
		perm, err := client.Permissions.Define().
			Object(complexObj).
			Action(complexAct).
			Create(ctx)
		require.NoError(t, err)
		
		err = client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
		require.NoError(t, err)
		defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)
		
		matchObj := "events:test" + uniqueSuffix + "/users:123"
		matchAct := "attendance:math:edit"
		
		allowed, err := client.Authz.Check().User(userID).Object(matchObj).Action(matchAct).Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed, "Should match complex patterns")
		
		allowed, err = client.Authz.Check().User(userID).Object(matchObj).Action("attendance:math:view").Allowed(ctx)
		require.NoError(t, err)
		assert.False(t, allowed, "Should NOT match different sub-action")
	})

	t.Run("Wildcard Matching Engine Check", func(t *testing.T) {
		// 1. Global Wildcard: * matches everything
		globalPerm, err := client.Permissions.Define().
			Object(goauth.ObjectWildcard()).
			Action(goauth.ActionWildcard()).
			Create(ctx)
		require.NoError(t, err)
		
		err = client.Permissions.GiveDirect(ctx, userID, globalPerm.ID, nil)
		require.NoError(t, err)
		
		allowed, _ := client.Authz.Check().User(userID).Object("any:thing").Action("any:action").Allowed(ctx)
		assert.True(t, allowed, "Global wildcard should match anything")
		
		client.Permissions.TakeDirect(ctx, userID, globalPerm.ID, nil)

		// 2. Recursive Object Wildcard: doc:**
		recursiveObj, _ := goauth.Object("doc", "rec"+uniqueSuffix).All().Build()
		recursivePerm, err := client.Permissions.Define().
			Object(recursiveObj).
			Action("read").
			Create(ctx)
		require.NoError(t, err)
		
		client.Permissions.GiveDirect(ctx, userID, recursivePerm.ID, nil)
		defer client.Permissions.TakeDirect(ctx, userID, recursivePerm.ID, nil)

		allowed, _ = client.Authz.Check().
			User(userID).
			Object("doc:rec" + uniqueSuffix + "/page:1/line:10").
			Action("read").
			Allowed(ctx)
		assert.True(t, allowed, "Recursive object wildcard should match deep hierarchy")

		// 3. Action Prefix Wildcard: file:*
		actionPrefixObj := "file:test" + uniqueSuffix
		prefixAct, _ := goauth.Action("file").Any().Build()
		prefixPerm, err := client.Permissions.Define().
			Object(actionPrefixObj).
			Action(prefixAct).
			Create(ctx)
		require.NoError(t, err)
		
		client.Permissions.GiveDirect(ctx, userID, prefixPerm.ID, nil)
		defer client.Permissions.TakeDirect(ctx, userID, prefixPerm.ID, nil)

		allowed, _ = client.Authz.Check().User(userID).Object(actionPrefixObj).Action("file:read").Allowed(ctx)
		assert.True(t, allowed, "Action wildcard should match sub-actions")
		
		allowed, _ = client.Authz.Check().User(userID).Object(actionPrefixObj).Action("file:write").Allowed(ctx)
		assert.True(t, allowed, "Action wildcard should match any sub-action")
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