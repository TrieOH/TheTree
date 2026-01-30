package testing

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/scopes"
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testScopes(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("scopes@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("ScopeProject")

	projectID := user.projectID

	var scopeID string
	t.Run("CreateScope", func(t *testing.T) {
		authClient := user.WithT(t)
		val := authClient.POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name":        "events",
				"external_id": "event1",
			}).
			Expect(http.StatusCreated).
			HasMessage("Scope Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &scopeID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "events",
			"external_id": "event1",
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GetScope", func(t *testing.T) {
		authClient := user.WithT(t)
		val := authClient.GET("/projects/" + projectID + "/scopes/" + scopeID).
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          AsString{Value: scopeID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "events",
			"external_id": "event1",
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateScopeNoName", func(t *testing.T) {
		user.WithT(t).POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"external_id": "event1",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.ScopeEmptyName).
			HasMessage("Scope name must not be empty")
	})

	t.Run("CreateScopeExternalIDAlreadyInName", func(t *testing.T) {
		user.WithT(t).POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name":        "events",
				"external_id": "event1",
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.ScopeDuplicateNameAndExternalID).
			HasMessage("two scopes can't have the same name and external_id")
	})

	t.Run("CreateScopeExistingNameNoID", func(t *testing.T) {
		authClient := user.WithT(t)
		val := authClient.POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name": "events",
			}).
			Expect(http.StatusCreated).
			HasMessage("Scope Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          AnyUUID{},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "events",
			"external_id": nil,
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateScopeExistingNameAndNewID", func(t *testing.T) {
		authClient := user.WithT(t)
		val := authClient.POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name":        "events",
				"external_id": "event2",
			}).
			Expect(http.StatusCreated).
			HasMessage("Scope Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          AnyUUID{},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "events",
			"external_id": "event2",
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GetAllProjectScopes", func(t *testing.T) {
		authClient := user.WithT(t)
		val := authClient.GET("/projects/" + projectID + "/scopes").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AnyUUID{},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "events",
				"external_id": "event1",
				"type":        "project_scope",
				"created_at":  AnyDate{},
			},
			map[string]interface{}{
				"id":          AnyUUID{},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "events",
				"external_id": nil,
				"type":        "project_scope",
				"created_at":  AnyDate{},
			},
			map[string]interface{}{
				"id":          AnyUUID{},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "events",
				"external_id": "event2",
				"type":        "project_scope",
				"created_at":  AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateGlobalScopeError", func(t *testing.T) {
		queries := sqlc.New(suite.DB)
		ctx := context.Background()
		_, err := queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  nil,
			Name:       nil,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsUniqueViolation(err))
		assert.Equal(t, "ERROR: duplicate key value violates unique constraint \"scopes_one_global\" (SQLSTATE 23505)", err.Error())

		pid, err := uuid.Parse(projectID)
		assert.NoError(t, err)

		nameStr := "global"
		externalIDStr := "global_id"

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  &pid,
			Name:       nil,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  nil,
			Name:       &nameStr,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  nil,
			Name:       nil,
			ExternalID: &externalIDStr,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  &pid,
			Name:       &nameStr,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  &pid,
			Name:       &nameStr,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  nil,
			Name:       &nameStr,
			ExternalID: &externalIDStr,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeGlobal),
			ProjectID:  &pid,
			Name:       &nameStr,
			ExternalID: &externalIDStr,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())
	})

	t.Run("CreateProjectRootScopeError", func(t *testing.T) {
		queries := sqlc.New(suite.DB)
		ctx := context.Background()

		pid, _ := uuid.Parse(projectID)
		nameStr := "global"
		externalIDStr := "global_id"

		_, err := queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeProjectRoot),
			ProjectID:  &pid,
			Name:       nil,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsUniqueViolation(err))
		assert.Equal(t, "ERROR: duplicate key value violates unique constraint \"scopes_one_project_root_per_project\" (SQLSTATE 23505)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeProjectRoot),
			ProjectID:  &pid,
			Name:       &nameStr,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeProjectRoot),
			ProjectID:  &pid,
			Name:       nil,
			ExternalID: &externalIDStr,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())

		_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       string(scopes.ScopeTypeProjectRoot),
			ProjectID:  &pid,
			Name:       &nameStr,
			ExternalID: &externalIDStr,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())
	})

	t.Run("CheckProjectRootScope", func(t *testing.T) {
		queries := sqlc.New(suite.DB)
		ctx := context.Background()
		pid, _ := uuid.Parse(projectID)
		rootScope, err := queries.GetRootByProjectID(ctx, &pid)
		assert.NoError(t, err)
		assert.EqualValues(t, &pid, rootScope.ProjectID)
		assert.Nil(t, rootScope.Name)
		assert.Nil(t, rootScope.ExternalID)
		assert.EqualValues(t, scopes.ScopeTypeProjectRoot, rootScope.Type)
	})

	t.Run("CreateInvalidScopeType", func(t *testing.T) {
		queries := sqlc.New(suite.DB)
		ctx := context.Background()

		pid, _ := uuid.Parse(projectID)

		_, err := queries.CreateScope(ctx, sqlc.CreateScopeParams{
			Type:       "invalid",
			ProjectID:  &pid,
			Name:       nil,
			ExternalID: nil,
		})

		assert.Error(t, err)
		assert.True(t, apierr.IsCheckViolation(err))
		assert.Equal(t, "ERROR: new row for relation \"scopes\" violates check constraint \"scope_shape_check\" (SQLSTATE 23514)", err.Error())
	})
}
