package goauth_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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
		uniqueSuffix := uuid.New().String()[:8]
		// 1. Global Wildcard: * matches everything
		// Use a unique action to avoid collision with existing permissions
		wildcardAction := "wildcard:action:" + uniqueSuffix
		globalPerm, err := client.Permissions.Define().
			Object(goauth.ObjectWildcard()).
			Action(wildcardAction).
			Create(ctx)
		require.NoError(t, err)

		err = client.Permissions.GiveDirect(ctx, userID, globalPerm.ID, nil)
		require.NoError(t, err)

		allowed, _ := client.Authz.Check().User(userID).Object("any:thing").Action(wildcardAction).Allowed(ctx)
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

	t.Run("Condition Builder Check", func(t *testing.T) {
		uniqueSuffix := strings.ReplaceAll(uuid.New().String()[:8], "-", "_")
		condObj := "doc:cond_" + uniqueSuffix

		// Define condition: resource.status == 'published' AND (resource.requester_id == resource.owner_id OR subject.role == 'admin')
		c := goauth.NewCondition()
		cond := c.And(
			c.Path("resource.status").Eq("published"),
			c.Or(
				c.Path("resource.requester_id").RefEq("resource.owner_id"),
				c.Path("subject.role").Eq("admin"),
			),
		)

		perm, err := client.Permissions.Define().
			Object(condObj).
			Action("read").
			Conditions(cond).
			Create(ctx)
		require.NoError(t, err)

		err = client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
		require.NoError(t, err)
		defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

		// 1. Match: status=published, owner_id=userID
		allowed, err := client.Authz.Check().
			User(userID).
			Object(condObj).
			Action("read").
			WithResource(map[string]any{
				"status":       "published",
				"owner_id":     userID.String(),
				"requester_id": userID.String(),
			}).
			Allowed(ctx)
		require.NoError(t, err)
		assert.True(t, allowed, "Should be allowed when owner and status match")

		// 2. No Match: status=draft, owner_id=userID
		allowed, err = client.Authz.Check().
			User(userID).
			Object(condObj).
			Action("read").
			WithResource(map[string]any{
				"status":       "draft",
				"owner_id":     userID.String(),
				"requester_id": userID.String(),
			}).
			Allowed(ctx)
		require.NoError(t, err)
		assert.False(t, allowed, "Should be denied when status is draft")

		// 3. No Match: status=published, owner_id=different
		allowed, err = client.Authz.Check().
			User(userID).
			Object(condObj).
			Action("read").
			WithResource(map[string]any{
				"status":       "published",
				"owner_id":     uuid.New().String(),
				"requester_id": userID.String(),
			}).
			Allowed(ctx)
		require.NoError(t, err)
		assert.False(t, allowed, "Should be denied when not owner and not admin")
	})

	t.Run("Simple Condition - Existence", func(t *testing.T) {
		uniqueSuffix := strings.ReplaceAll(uuid.New().String()[:8], "-", "_")
		obj := "doc:simple_" + uniqueSuffix
		c := goauth.NewCondition()

		perm, _ := client.Permissions.Define().
			Object(obj).
			Action("read").
			Conditions(c.Path("resource.metadata").Exists()).
			Create(ctx)

		_ = client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
		defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

		allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("read").
			WithResource(map[string]any{"metadata": map[string]any{"key": "val"}}).Allowed(ctx)
		assert.True(t, allowed)

		allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
			WithResource(map[string]any{"other": "thing"}).Allowed(ctx)
		assert.False(t, allowed)
	})

	t.Run("Complex Condition - Temporal and Logical", func(t *testing.T) {
		uniqueSuffix := strings.ReplaceAll(uuid.New().String()[:8], "-", "_")
		obj := "doc:complex_" + uniqueSuffix
		c := goauth.NewCondition()

		// (status == 'active' AND grace_around 1h) OR role == 'manager'
		cond := c.Or(
			c.And(
				c.Path("resource.status").Eq("active"),
				c.Field("resource.valid_at").GraceAround("1h"),
			),
			c.Path("subject.role").Eq("manager"),
		)

		perm, _ := client.Permissions.Define().
			Object(obj).
			Action("write").
			Conditions(cond).
			Create(ctx)

		_ = client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
		defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

		now := time.Now().UTC()

		// Match: active + within grace
		allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("write").
			WithResource(map[string]any{
				"status":   "active",
				"valid_at": now.Add(30 * time.Minute).Format(time.RFC3339),
			}).Allowed(ctx)
		assert.True(t, allowed)

		// No match: active but outside grace
		allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("write").
			WithResource(map[string]any{
				"status":   "active",
				"valid_at": now.Add(2 * time.Hour).Format(time.RFC3339),
			}).Allowed(ctx)
		assert.False(t, allowed)
	})

	t.Run("Very Complex Condition", func(t *testing.T) {
		uniqueSuffix := strings.ReplaceAll(uuid.New().String()[:8], "-", "_")
		obj := "doc:vcomplex_" + uniqueSuffix
		c := goauth.NewCondition()

		/*
		  AND(
		    priority >= 10,
		    OR(
		      requester_id == owner_id,
		      groups containsAny ['admin', 'editor']
		    ),
		    NOT(locked == true),
		    expires_at grace_before 24h
		  )
		*/
		cond := c.And(
			c.Path("resource.priority").Gte(10),
			c.Or(
				c.Path("resource.requester_id").RefEq("resource.owner_id"),
				c.Path("resource.groups").ContainsAny([]string{"admin", "editor"}),
			),
			c.Not(c.Path("resource.locked").Eq(true)),
			c.Field("resource.expires_at").GraceBefore("24h"),
		)

		perm, err := client.Permissions.Define().
			Object(obj).
			Action("delete").
			Conditions(cond).
			Create(ctx)
		require.NoError(t, err)

		_ = client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
		defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

		now := time.Now().UTC()
		future := now.Add(12 * time.Hour).Format(time.RFC3339)

		// 1. Full Match
		allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("delete").
			WithResource(map[string]any{
				"priority":     15,
				"owner_id":     "user1",
				"requester_id": "user1",
				"locked":       false,
				"expires_at":   future,
			}).Allowed(ctx)
		assert.True(t, allowed)

		// 2. Fail on priority
		allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("delete").
			WithResource(map[string]any{
				"priority":     5,
				"owner_id":     "user1",
				"requester_id": "user1",
				"locked":       false,
				"expires_at":   future,
			}).Allowed(ctx)
		assert.False(t, allowed)

		// 3. Match via groups
		allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("delete").
			WithResource(map[string]any{
				"priority":     20,
				"owner_id":     "user1",
				"requester_id": "user2", // not owner
				"groups":       []string{"editor", "viewer"},
				"locked":       false,
				"expires_at":   future,
			}).Allowed(ctx)
		assert.True(t, allowed)

		// 4. Fail on locked
		allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("delete").
			WithResource(map[string]any{
				"priority":     15,
				"owner_id":     "user1",
				"requester_id": "user1",
				"locked":       true,
				"expires_at":   future,
			}).Allowed(ctx)
		assert.False(t, allowed)
	})

	t.Run("Grace Operators Check", func(t *testing.T) {
		uniqueSuffix := strings.ReplaceAll(uuid.New().String()[:8], "-", "_")
		c := goauth.NewCondition()

		t.Run("GraceBefore", func(t *testing.T) {
			obj := "doc:gb_" + uniqueSuffix
			perm, _ := client.Permissions.Define().Object(obj).Action("read").
				Conditions(c.Field("resource.ts").GraceBefore("1h")).Create(ctx)
			client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
			defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

			now := time.Now().UTC()
			// Valid: [ts - 1h, ts] -> so ts must be in [now, now + 1h]
			
			// Valid: ts is 30m in future, so now is 30m before ts (within 1h)
			allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(30 * time.Minute).Format(time.RFC3339)}).Allowed(ctx)
			assert.True(t, allowed)

			// Invalid: ts is 2h in future, so now is 2h before ts (outside 1h)
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(2 * time.Hour).Format(time.RFC3339)}).Allowed(ctx)
			assert.False(t, allowed)

			// Invalid: ts is in past, so now is AFTER ts
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(-10 * time.Minute).Format(time.RFC3339)}).Allowed(ctx)
			assert.False(t, allowed)
		})

		t.Run("GraceAfter", func(t *testing.T) {
			obj := "doc:ga_" + uniqueSuffix
			perm, _ := client.Permissions.Define().Object(obj).Action("read").
				Conditions(c.Field("resource.ts").GraceAfter("1h")).Create(ctx)
			client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
			defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

			now := time.Now().UTC()
			// Valid: [ts, ts + 1h] -> so ts must be in [now - 1h, now]

			// Valid: ts is 30m in past
			allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(-30 * time.Minute).Format(time.RFC3339)}).Allowed(ctx)
			assert.True(t, allowed)

			// Invalid: ts is 2h in past
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(-2 * time.Hour).Format(time.RFC3339)}).Allowed(ctx)
			assert.False(t, allowed)
		})

		t.Run("GraceAround", func(t *testing.T) {
			obj := "doc:gar_" + uniqueSuffix
			perm, _ := client.Permissions.Define().Object(obj).Action("read").
				Conditions(c.Field("resource.ts").GraceAround("1h")).Create(ctx)
			client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
			defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

			now := time.Now().UTC()
			// Valid: [ts - 1h, ts + 1h]

			// Valid: 30m before
			allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(30 * time.Minute).Format(time.RFC3339)}).Allowed(ctx)
			assert.True(t, allowed)

			// Valid: 30m after
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(-30 * time.Minute).Format(time.RFC3339)}).Allowed(ctx)
			assert.True(t, allowed)

			// Invalid: 2h away
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{"ts": now.Add(2 * time.Hour).Format(time.RFC3339)}).Allowed(ctx)
			assert.False(t, allowed)
		})

		t.Run("GraceDuration", func(t *testing.T) {
			obj := "doc:gd_" + uniqueSuffix
			perm, _ := client.Permissions.Define().Object(obj).Action("read").
				Conditions(c.Fields("resource.start", "resource.end").GraceDuration("15m")).Create(ctx)
			client.Permissions.GiveDirect(ctx, userID, perm.ID, nil)
			defer client.Permissions.TakeDirect(ctx, userID, perm.ID, nil)

			now := time.Now().UTC()
			// Valid: [start - 15m, end + 15m]

			// Valid: Right in the middle
			allowed, _ := client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{
					"start": now.Add(-1 * time.Hour).Format(time.RFC3339),
					"end":   now.Add(1 * time.Hour).Format(time.RFC3339),
				}).Allowed(ctx)
			assert.True(t, allowed)

			// Valid: 5m before start
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{
					"start": now.Add(5 * time.Minute).Format(time.RFC3339),
					"end":   now.Add(1 * time.Hour).Format(time.RFC3339),
				}).Allowed(ctx)
			assert.True(t, allowed)

			// Invalid: 30m before start
			allowed, _ = client.Authz.Check().User(userID).Object(obj).Action("read").
				WithResource(map[string]any{
					"start": now.Add(30 * time.Minute).Format(time.RFC3339),
					"end":   now.Add(1 * time.Hour).Format(time.RFC3339),
				}).Allowed(ctx)
			assert.False(t, allowed)
		})
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
